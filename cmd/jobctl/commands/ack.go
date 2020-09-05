package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
	"github.com/jeffrom/job-manager/pkg/job"
)

type ackCmd struct {
	*cobra.Command
}

func newAckCmd(cfg *jobclient.Config) *ackCmd {
	c := &ackCmd{
		Command: &cobra.Command{
			Use: "ack",
		},
	}

	return c
}

func (c *ackCmd) Cmd() *cobra.Command { return c.Command }
func (c *ackCmd) Execute(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
	client := clientFromContext(ctx)
	status := job.StatusComplete
	return client.AckJob(ctx, args[0], status)
}
