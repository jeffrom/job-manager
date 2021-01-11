package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
)

type getCmd struct {
	*cobra.Command
	gqCmd *getQueueCmd
}

func newGetCmd(cfg *client.Config) *getCmd {
	gqCmd := newGetQueueCmd(cfg)
	cmd := gqCmd.Cmd()
	cmd.Use = "get"
	cmd.Aliases = []string{"g"}
	c := &getCmd{
		gqCmd:   gqCmd,
		Command: cmd,
	}

	addCommands(cfg, c,
		newGetQueueCmd(cfg),
	)
	return c
}

func (c *getCmd) Cmd() *cobra.Command { return c.Command }
func (c *getCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	return c.gqCmd.Execute(ctx, cfg, cmd, args)
}
