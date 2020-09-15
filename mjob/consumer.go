package mjob

import (
	"context"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/pkg/resource"
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
		wrk := newWorker(c.cfg, ch, c.resultC)
		workers[i] = wrk
		wrk.start(ctx)
	}
	c.workers = workers

	curr := make([]*resource.Job, c.cfg.Concurrency+c.cfg.Backpressure)
	for {
		select {
		case <-c.stop:
			// do the shutdown sequence
			return nil
		case <-ctx.Done():
			// trigger the shutdown sequence
			return c.Stop()
		}

		n := (c.cfg.Concurrency + c.cfg.Backpressure) - (int(c.running) + len(curr))
		var jobs *resource.Jobs
		if n > 0 {
			var err error
			jobs, err = c.client.DequeueJobsOpts(ctx, n, c.cfg.DequeueOpts)
			if err != nil {
				// TODO backoff
				return err
			}
		}

		remaining, err := c.processJobs(ctx, append(curr, jobs.Jobs...))
		if err != nil {
			// TODO handle this? should processJobs ever error? only if a
			// worker doesn't finish in time.
			return err
		}
		curr = remaining
		return nil
	}
	return nil
}

// TODO processJobs should return when there is at least one completed
// job, so that subsequent jobs can start on that worker, but it should
// start as many jobs as possible first.
func (c *Consumer) processJobs(ctx context.Context, jobs []*resource.Job) ([]*resource.Job, error) {
	currJobs := jobs
	for {
		if len(currJobs) == 0 {
			break
		}
		jb := currJobs[0]

		if c.startJob(ctx, jb) {
			currJobs = currJobs[1:]
		}
	}

	// TODO wait up to JobDuration, then throw it in a penalty box where
	// concurrency is decreased by one until the worker completes? To start
	// maybe just stop the consumer.
	// finished := make(map[string]bool)
	for {
		select {
		case <-c.resultC:
		default:
			break
		}
	}
	return jobs, nil
}

func (c *Consumer) startJob(ctx context.Context, jb *resource.Job) bool {
	started := false
	for _, wrk := range c.workers {
		select {
		case wrk.in <- jb:
			started = true
			break
		default:
		}
	}
	return started
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
