package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
)

type listCmd struct {
	*cobra.Command
	lqCmd *listQueuesCmd
}

func newListCmd(cfg *jobclient.Config) *listCmd {
	lqCmd := newListQueuesCmd(cfg)
	cmd := lqCmd.Cmd()
	cmd.Use = "list"
	cmd.Aliases = []string{"ls", "get"}
	c := &listCmd{
		lqCmd:   lqCmd,
		Command: cmd,
	}

	addCommands(cfg, c,
		newListQueuesCmd(cfg),
	)
	return c
}

func (c *listCmd) Cmd() *cobra.Command { return c.Command }
func (c *listCmd) Execute(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
	return c.lqCmd.Execute(ctx, cfg, cmd, args)
}