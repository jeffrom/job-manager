package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
)

type showCmd struct {
	*cobra.Command
	showJobCmd *showJobCmd
}

func newShowCmd(cfg *client.Config) *showCmd {
	showJobCmd := newShowJobCmd(cfg)
	cmd := showJobCmd.Cmd()
	cmd.Use = "show"
	cmd.Aliases = []string{"describe", "desc"}
	c := &showCmd{
		showJobCmd: showJobCmd,
		Command:    cmd,
	}

	addCommands(cfg, c,
		newShowQueueCmd(cfg),
	)
	return c
}

func (c *showCmd) Cmd() *cobra.Command { return c.Command }
func (c *showCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	return c.showJobCmd.Execute(ctx, cfg, cmd, args)
}
