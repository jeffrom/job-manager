// Package consumer implements a job executor.
package consumer

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/logger"
	"github.com/jeffrom/job-manager/mjob/resource"
)

// Runner can be implemented to execute jobs.
type Runner interface {
	// Run executes a job. It can return a nil *JobResult, indicating success.
	// Returning an error will cause the job to fail.
	Run(ctx context.Context, job *resource.Job) (*resource.JobResult, error)
}

// Consumer manages the dequeue, execution, and acknowledgement of jobs.
type Consumer struct {
	cfg     Config
	log     logger.Logger
	client  client.Interface
	runner  Runner
	resultC chan *resource.JobResult
	workers []*worker
	running int32
	stop    chan struct{}

	activeJobs []string
	mu         sync.Mutex
}

// Provider provides configuration to Consumers. It can be passed as an
// argument to New.
type Provider func(c *Consumer) *Consumer

// New creates a new instance of Consumer.
func New(client client.Interface, runner Runner, providers ...Provider) *Consumer {
	c := &Consumer{
		client: client,
		runner: runner,
		cfg:    defaultConfig,
		log:    &logger.Default{},
		stop:   make(chan struct{}, 1),
	}

	for _, provider := range providers {
		c = provider(c)
	}

	size := c.cfg.Concurrency
	if size == 0 {
		size = 1
	}
	c.resultC = make(chan *resource.JobResult, size)
	return c
}

// WithConfig allows custom configuration to be provided to new Consumers.
func WithConfig(cfg Config) Provider {
	return func(c *Consumer) *Consumer {
		c.cfg = cfg
		return c
	}
}

// WithLogger allows a custom logger to be provided to new Consumers.
func WithLogger(logger logger.Logger) Provider {
	return func(c *Consumer) *Consumer {
		c.log = logger
		return c
	}
}

// WithQueue adds an additional target queue for new Consumers. If no queues
// are provided, all queues are dequeued from. Multiple queues can be consumed.
func WithQueue(queue string) Provider {
	return func(c *Consumer) *Consumer {
		c.cfg.DequeueOpts.Queues = append(c.cfg.DequeueOpts.Queues, queue)
		return c
	}
}

// Run consumes jobs until Stop is called. After Stop is called, any currently
// running jobs will continue until completion.
func (c *Consumer) Run(ctx context.Context) error {
	n := c.cfg.Concurrency
	if n == 0 {
		n = 1
	}
	workers := make([]*worker, n)
	for i := 0; i < n; i++ {
		ch := make(chan *resource.Job)
		wrk := newWorker(c.cfg, c.log, c.runner, ch, c.resultC)
		workers[i] = wrk
		go wrk.start(ctx)
	}
	c.workers = workers

	size := c.cfg.Concurrency
	if size == 0 {
		size = 1
	}
	curr := make([]*resource.Job, 0, size)
Loop:
	for {
		select {
		case <-c.stop:
			// fmt.Println("<- c.stop")
			break Loop
		case <-ctx.Done():
			// fmt.Println("<- ctx.Done")
			// trigger the shutdown sequence
			break Loop
		default:
		}

		// fmt.Println("counts", c.cfg.Concurrency, c.cfg.Backpressure, c.running, len(curr))
		n := (maxInt(c.cfg.Concurrency, 1)) - (int(atomic.LoadInt32(&c.running)) + len(curr))
		// fmt.Println("will get", n, "jobs")
		jobs := &resource.Jobs{}
		if n > 0 {
			var err error
			jobs, err = c.client.DequeueJobsOpts(ctx, n, c.cfg.DequeueOpts)
			if err != nil {
				// TODO backoff
				c.log.Log(ctx, &logger.Event{
					Level:   "error",
					Error:   err,
					Message: "dequeue failed",
				})

				sleep(ctx, 2*time.Second)
				continue
			}
		}
		dequeuedJobs := len(jobs.Jobs)

		// fmt.Println("consumer got", dequeuedJobs, "jobs with limit", n)
		remaining, err := c.processJobs(ctx, append(curr, jobs.Jobs...))
		// fmt.Println("consumer processed jobs.", len(remaining), "jobs remain", "err:", err)
		if err != nil {
			// TODO handle this? should processJobs ever error? only if a
			// worker doesn't finish in time.
			c.log.Log(ctx, &logger.Event{
				Level:   "error",
				Error:   err,
				Message: "processJobs failed",
			})
			continue
		}
		curr = remaining
		if dequeuedJobs == 0 && n != 0 {
			sleep(ctx, 2*time.Second)
		}
	}

	currRunning := atomic.LoadInt32(&c.running)
	if currRunning == 0 {
		c.log.Log(ctx, &logger.Event{
			Level:   "info",
			Message: "skipping shutdown sequence as there are 0 jobs running",
		})
		return nil
	}

	c.log.Log(ctx, &logger.Event{
		Level:   "info",
		Message: "shutdown sequence beginning",
		Data:    map[string]int32{"running": currRunning},
	})

	shutdownCtx, shutdownDone := context.WithDeadline(context.Background(), time.Now().Add(c.cfg.ShutdownTimeout))
	defer shutdownDone()
ShutdownLoop:
	// cleanup
	for {
		select {
		case <-shutdownCtx.Done():
			n := atomic.LoadInt32(&c.running)
			finalCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
			defer cancel()
			for _, jobID := range c.getActive() {
				if err := c.client.AckJob(finalCtx, jobID, resource.StatusFailed); err != nil {
					c.log.Log(shutdownCtx, &logger.Event{
						Level:   "error",
						Error:   err,
						Message: "ackJob failed",
					})
				}
			}
			c.log.Log(context.Background(), &logger.Event{
				Level:   "error",
				Message: fmt.Sprintf("shut down with %d jobs still in progress", n),
				Data:    map[string]int32{"running": n},
			})
			break ShutdownLoop
		case res := <-c.resultC:
			n := atomic.AddInt32(&c.running, -1)
			c.removeActive(res.JobID)
			if err := c.client.AckJobOpts(shutdownCtx, res.JobID, *res.Status, client.AckJobOpts{Data: res.Data}); err != nil {
				c.log.Log(shutdownCtx, &logger.Event{
					Level:   "error",
					Error:   err,
					Message: "ackJob failed",
				})
			}
			if n == 0 {
				break ShutdownLoop
			}
		}
	}
	for _, wrk := range c.workers {
		close(wrk.in)
	}
	c.log.Log(shutdownCtx, &logger.Event{
		Level:   "info",
		Message: "shutdown sequence complete",
	})
	return nil
}

// TODO processJobs should return when there is at least one completed
// job, so that subsequent jobs can start on that worker, but it should
// start as many jobs as possible first.
func (c *Consumer) processJobs(ctx context.Context, jobs []*resource.Job) ([]*resource.Job, error) {
	currJobs := make([]*resource.Job, len(jobs))
	copy(currJobs, jobs)
	for {
		if len(currJobs) == 0 {
			break
		}
		jb := currJobs[0]

		if c.startJob(ctx, jb) {
			currJobs = currJobs[1:]
		} else {
			break
		}
	}

Loop:
	for {
		select {
		case res := <-c.resultC:
			atomic.AddInt32(&c.running, -1)
			if res == nil || res.JobID == "" {
				panic("consumer: job id was empty")
			}
			// TODO we should make an effort to ack all jobs, but it is always
			// possible for a job to run twice. It runs in its own context.
			if err := c.ackJob(ctx, res.JobID, *res.Status, res); err != nil {
				return currJobs, err
			}
			c.removeActive(res.JobID)
			for i, jb := range currJobs {
				if jb.ID == res.JobID {
					// remove job from processing list
					currJobs = append(currJobs[:i], currJobs[i+1:]...)
				}
			}
			// fmt.Printf("processJobs: got a result: %+v\n", res)
		default:
			// fmt.Println("processJobs: got no results")
			break Loop
		}
	}

	return currJobs, nil
}

func (c *Consumer) ackJob(ctx context.Context, jobID string, status resource.Status, res *resource.JobResult) error {
	if cerr := ctx.Err(); cerr != nil {
		nctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(c.cfg.ShutdownTimeout))
		defer cancel()
		ctx = nctx
	}
	return c.client.AckJobOpts(ctx, res.JobID, *res.Status, client.AckJobOpts{Data: res.Data})
}

func (c *Consumer) startJob(ctx context.Context, jb *resource.Job) bool {
	for _, wrk := range c.workers {
		select {
		case wrk.in <- jb:
			atomic.AddInt32(&c.running, 1)
			c.addActive(jb.ID)
			return true
		default:
		}
	}
	return false
}

// Stop initiates the consumer's shutdown sequence.
func (c *Consumer) Stop() {
	c.stop <- struct{}{}
}

func (c *Consumer) addActive(ids ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, id := range ids {
		if _, ok := inStrings(id, c.activeJobs); ok {
			continue
		}
		c.activeJobs = append(c.activeJobs, id)
	}
}

func (c *Consumer) removeActive(ids ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, id := range ids {
		if i, ok := inStrings(id, c.activeJobs); ok {
			c.activeJobs = append(c.activeJobs[:i], c.activeJobs[i+1:]...)
		}
	}
}

func (c *Consumer) getActive() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	cp := make([]string, len(c.activeJobs))
	copy(cp, c.activeJobs)
	return cp
}

func inStrings(a string, b []string) (int, bool) {
	for i, s := range b {
		if a == s {
			return i, true
		}
	}
	return -1, false
}

func sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
