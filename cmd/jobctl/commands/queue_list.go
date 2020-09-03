package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
)

type queueListOpts struct {
}

type queueListCmd struct {
	*cobra.Command
	opts *queueListOpts
}

func newQueueListCmd(cfg *jobclient.Config) *queueListCmd {
	opts := &queueListOpts{}
	c := &queueListCmd{
		opts: opts,
		Command: &cobra.Command{
			Use:     "list",
			Aliases: []string{"ls"},
			Args:    cobra.NoArgs,
		},
	}

	return c
}

func (c *queueListCmd) Cmd() *cobra.Command { return c.Command }
func (c *queueListCmd) Execute(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
	return nil
}
