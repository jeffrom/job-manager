package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
)

type queueListCmd struct {
	*cobra.Command
}

func (c *queueListCmd) Cmd() *cobra.Command { return c.Command }
func (c *queueListCmd) Execute(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
	return nil
}

func newQueueListCmd(cfg *jobclient.Config) *queueListCmd {
	c := &queueListCmd{
		Command: &cobra.Command{
			Use:     "list",
			Aliases: []string{"ls"},
			Args:    cobra.MaximumNArgs(1),
		},
	}

	return c
}
