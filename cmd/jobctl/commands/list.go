package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
)

type listCmd struct {
	*cobra.Command
}

func newListCmd(cfg *jobclient.Config) *listCmd {
	c := &listCmd{
		Command: &cobra.Command{
			Use:     "list",
			Aliases: []string{"ls", "get"},
		},
	}

	addCommands(cfg, c,
		newListQueuesCmd(cfg),
	)
	return c
}

func (c *listCmd) Cmd() *cobra.Command { return c.Command }
func (c *listCmd) Execute(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
	return runListQueues(ctx, cfg, cmd, args)
}
