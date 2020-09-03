package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
)

type queueCmd struct {
	*cobra.Command
}

func (c *queueCmd) Cmd() *cobra.Command { return c.Command }
func (c *queueCmd) Execute(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
	return usageCmd(ctx, cfg, cmd, args)
}

func newQueueCmd(cfg *jobclient.Config) *queueCmd {
	c := &queueCmd{
		Command: &cobra.Command{
			Use:     "queue",
			Aliases: []string{"q"},
		},
	}

	addCommands(cfg, c,
		newQueueListCmd(cfg),
		newQueueSaveCmd(cfg),
	)
	return c
}
