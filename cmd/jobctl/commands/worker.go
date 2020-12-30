package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob"
	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/resource"
)

type workerOpts struct {
	// client.workerOpts
	// data string
	// claims []string
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

	// flags := c.Command.Flags()
	// flags.StringArrayVarP(&opts.claims, "claim", "c", nil, "worker with claims")

	return c
}

func (c *workerCmd) Cmd() *cobra.Command { return c.Command }
func (c *workerCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	cl := clientFromContext(ctx)
	queues := args
	consumer := mjob.NewConsumer(cl, &runner{}, mjob.ConsumerWithConfig(mjob.ConsumerConfig{
		Concurrency: 1,
		DequeueOpts: client.DequeueOpts{
			Queues: queues,
		},
	}))

	return consumer.Run(ctx)
}

type runner struct{}

func (r *runner) Run(ctx context.Context, job *resource.Job) (*resource.JobResult, error) {
	fmt.Printf("job: %+v\n", job)
	return nil, nil
}
