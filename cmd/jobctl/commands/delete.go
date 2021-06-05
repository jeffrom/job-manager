package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
)

type deleteCmd struct {
	*cobra.Command
	deleteQueueCmd *deleteQueueCmd
}

func newDeleteCmd(cfg *client.Config) *deleteCmd {
	deletec := newDeleteQueueCmd(cfg)
	cmd := deletec.Cmd()
	cmd.Use = "delete"
	cmd.Aliases = []string{"del"}
	cmd.Args = cobra.MaximumNArgs(1)
	c := &deleteCmd{
		deleteQueueCmd: deletec,
		Command:        cmd,
	}

	addCommands(cfg, c,
		newDeleteQueueCmd(cfg),
	)
	return c
}

func (c *deleteCmd) Cmd() *cobra.Command { return c.Command }
func (c *deleteCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
