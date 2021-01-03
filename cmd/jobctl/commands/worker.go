package commands

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob"
	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/label"
	"github.com/jeffrom/job-manager/mjob/resource"
)

type workerOpts struct {
	concurrency int
	claims      []string
	sleep       time.Duration
	failTimes   int
}

type workerCmd struct {
	*cobra.Command
	opts *workerOpts
}

func newWorkerCmd(cfg *client.Config) *workerCmd {
	opts := &workerOpts{}
	c := &workerCmd{
		Command: &cobra.Command{
			Use:  "worker QUEUE...",
			Args: cobra.MinimumNArgs(1),
			// Aliases: []string{"wrk"},
		},
		opts: opts,
	}

	flags := c.Command.Flags()
	flags.IntVarP(&opts.concurrency, "concurrency", "C", 1, "max concurrent jobs")
	flags.IntVar(&opts.failTimes, "fail-times", 0, "number of failures before success")
	flags.StringArrayVarP(&opts.claims, "claim", "c", nil, "claims for this worker")
	flags.DurationVar(&opts.sleep, "sleep", 0, "sleep before completion")

	return c
}

func (c *workerCmd) Cmd() *cobra.Command { return c.Command }
func (c *workerCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	var claims label.Claims
	if len(c.opts.claims) > 0 {
		var err error
		claims, err = label.ParseClaims(c.opts.claims)
		if err != nil {
			return err
		}
	}
	cl := clientFromContext(ctx)
	queues := args
	consumer := mjob.NewConsumer(cl, &runner{opts: c.opts}, mjob.ConsumerWithConfig(mjob.ConsumerConfig{
		Concurrency: c.opts.concurrency,
		DequeueOpts: client.DequeueOpts{
			Claims: claims,
			Queues: queues,
		},
	}))

	log.Print("Starting worker on queues: ", strings.Join(queues, ", "))
	return consumer.Run(ctx)
}

type runner struct {
	opts *workerOpts
}

func (r *runner) Run(ctx context.Context, job *resource.Job) (*resource.JobResult, error) {
	log.Printf("job %s: %+v", job.ID, job)
	if r.opts.sleep > 0 {
		time.Sleep(r.opts.sleep)
	}
	fmt.Println(r.opts.failTimes, job.Attempt)
	if ft := r.opts.failTimes; ft > 0 && job.Attempt <= ft {
		log.Printf("failing job %s on attempt %d", job.ID, job.Attempt)
		return nil, fmt.Errorf("failing job attempt %d", job.Attempt)
	}
	log.Printf("job %s complete", job.ID)
	return nil, nil
}
