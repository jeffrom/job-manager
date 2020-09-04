package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
)

type saveCmd struct {
	*cobra.Command
}

func (c *saveCmd) Cmd() *cobra.Command { return c.Command }
func (c *saveCmd) Execute(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
	return usageCmd(ctx, cfg, cmd, args)
}

func newSaveCmd(cfg *jobclient.Config) *saveCmd {
	c := &saveCmd{
		Command: &cobra.Command{
			Use: "save",
		},
	}

	addCommands(cfg, c,
		newSaveQueueCmd(cfg),
	)
	return c
}
