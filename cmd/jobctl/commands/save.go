package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
)

type saveCmd struct {
	*cobra.Command
	saveQueueCmd *saveQueueCmd
}

func newSaveCmd(cfg *client.Config) *saveCmd {
	savec := newSaveQueueCmd(cfg)
	cmd := savec.Cmd()
	cmd.Use = "save"
	cmd.Args = cobra.MaximumNArgs(1)
	c := &saveCmd{
		saveQueueCmd: savec,
		Command:      cmd,
	}

	addCommands(cfg, c,
		newSaveQueueCmd(cfg),
	)
	return c
}

func (c *saveCmd) Cmd() *cobra.Command { return c.Command }
func (c *saveCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Usage()
	}
	return c.saveQueueCmd.Execute(ctx, cfg, cmd, args)
}
