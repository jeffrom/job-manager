package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
)

type listCmd struct {
	*cobra.Command
}

func (c *listCmd) Cmd() *cobra.Command { return c.Command }
func (c *listCmd) Execute(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
	return usageCmd(ctx, cfg, cmd, args)
}

func newListCmd(cfg *jobclient.Config) *listCmd {
	c := &listCmd{
		Command: &cobra.Command{
			Use:     "list",
			Aliases: []string{"ls"},
		},
	}

	addCommands(cfg, c,
		newListQueuesCmd(cfg),
	)
	return c
}
