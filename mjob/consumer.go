package mjob

import (
	"context"
	"time"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/resource"
)

type Runner interface {
	Run(ctx context.Context, job *resource.Job) (*resource.JobResult, error)
}

type Consumer struct {
	cfg     ConsumerConfig
	client  client.Interface
	runner  Runner
	resultC chan *resource.JobResult
	workers []*consumerWorker
	running int32
	stop    chan struct{}
}

type ConsumerProviderFunc func(c *Consumer) *Consumer

func NewConsumer(client client.Interface, runner Runner, providers ...ConsumerProviderFunc) *Consumer {
	c := &Consumer{
		client: client,
		runner: runner,
		cfg:    defaultConsumerConfig,
		stop:   make(chan struct{}),
	}

	for _, provider := range providers {
		c = provider(c)
	}

	if c.cfg.Concurrency == 0 {
		panic("concurrency config was 0")
	}
	c.resultC = make(chan *resource.JobResult, c.cfg.Concurrency)
	if c.cfg.Backpressure == 0 {
		c.cfg.Backpressure = c.cfg.Concurrency / 2
	}
	return c
}

func ConsumerWithConfig(cfg ConsumerConfig) ConsumerProviderFunc {
	return func(c *Consumer) *Consumer {
		c.cfg = cfg
		return c
	}
}

// Run consumes jobs until Stop is called. After Stop is called, any currently
// running jobs will continue until completion.
func (c *Consumer) Run(ctx context.Context) error {
	n := c.cfg.Concurrency
	workers := make([]*consumerWorker, n)
	for i := 0; i < n; i++ {
		ch := make(chan *resource.Job)
		wrk := newWorker(c.cfg, c.runner, ch, c.resultC)
		workers[i] = wrk
		go wrk.start(ctx)
	}
	c.workers = workers

	curr := make([]*resource.Job, 0, c.cfg.Concurrency+c.cfg.Backpressure)
	for {
		select {
		case <-c.stop:
			// do the shutdown sequence
			return nil
		case <-ctx.Done():
			// trigger the shutdown sequence
			return c.Stop()
		default:
		}

		// fmt.Println("counts", c.cfg.Concurrency, c.cfg.Backpressure, c.running, len(curr))
		n := (c.cfg.Concurrency) - (int(c.running) + len(curr))
		// fmt.Println("will get", n, "jobs")
		jobs := &resource.Jobs{}
		if n > 0 {
			var err error
			jobs, err = c.client.DequeueJobsOpts(ctx, n, c.cfg.DequeueOpts)
			if err != nil {
				// TODO backoff
				return err
			}
		}

		njobs := len(jobs.Jobs)
		// fmt.Println("consumer got", njobs, "jobs")
		remaining, err := c.processJobs(ctx, append(curr, jobs.Jobs...))
		// fmt.Println("consumer processed jobs.", len(remaining), "jobs remain", "err:", err)
		if err != nil {
			// TODO handle this? should processJobs ever error? only if a
			// worker doesn't finish in time.
			return err
		}
		curr = remaining
		if njobs == 0 {
			time.Sleep(2 * time.Second)
		}
	}
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

	// TODO wait up to JobDuration, then throw it in a penalty box where
	// concurrency is decreased by one until the worker completes? To start
	// maybe just stop the consumer.
	// maybe just backoff on acks?
	// finished := make(map[string]bool)
Loop:
	for {
		select {
		case res := <-c.resultC:
			if res == nil || res.JobID == "" {
				panic("consumer: job id was empty")
			}
			// we should make an effort to ack all jobs, but it is always
			// possible for a job to run twice
			if err := c.client.AckJobOpts(ctx, res.JobID, *res.Status, client.AckJobOpts{Data: res.Data}); err != nil {
				return currJobs, err
			}
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

func (c *Consumer) startJob(ctx context.Context, jb *resource.Job) bool {
	for _, wrk := range c.workers {
		select {
		case wrk.in <- jb:
			return true
		default:
		}
	}
	return false
}

// Stop cancels any pending jobs, waits for any currently running jobs to
// complete, then stops the consumer.
func (c *Consumer) Stop() error {
	for _, wrk := range c.workers {
		close(wrk.in)
	}

	// TODO cancel pending jobs

	// TODO process remaining jobs
	return nil
}
