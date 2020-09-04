package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
)

type saveCmd struct {
	*cobra.Command
	saveQueueCmd *saveQueueCmd
}

func (c *saveCmd) Cmd() *cobra.Command { return c.Command }
func (c *saveCmd) Execute(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
	return c.saveQueueCmd.Execute(ctx, cfg, cmd, args)
}

func newSaveCmd(cfg *jobclient.Config) *saveCmd {
	savec := newSaveQueueCmd(cfg)
	cmd := savec.Cmd()
	cmd.Use = "save"
	c := &saveCmd{
		saveQueueCmd: savec,
		Command:      cmd,
	}

	addCommands(cfg, c,
		newSaveQueueCmd(cfg),
	)
	return c
}
