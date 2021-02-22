package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
)

type listCmd struct {
	*cobra.Command
	ljCmd *listJobsCmd
}

func newListCmd(cfg *client.Config) *listCmd {
	ljCmd := newListJobsCmd(cfg)
	cmd := ljCmd.Cmd()
	cmd.Use = "list"
	cmd.Aliases = []string{"ls"}
	c := &listCmd{
		ljCmd:   ljCmd,
		Command: cmd,
	}

	addCommands(cfg, c,
		newListQueuesCmd(cfg),
	)
	return c
}

func (c *listCmd) Cmd() *cobra.Command { return c.Command }
func (c *listCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	return c.ljCmd.Execute(ctx, cfg, cmd, args)
}
