package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/resource"
)

type ackCmd struct {
	*cobra.Command
}

func newAckCmd(cfg *client.Config) *ackCmd {
	c := &ackCmd{
		Command: &cobra.Command{
			Use: "ack",
		},
	}

	return c
}

func (c *ackCmd) Cmd() *cobra.Command { return c.Command }
func (c *ackCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	client := clientFromContext(ctx)
	status := resource.StatusComplete
	return client.AckJob(ctx, args[0], status)
}
