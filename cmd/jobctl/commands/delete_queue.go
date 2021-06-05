package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
)

type deleteQueueOpts struct {
}

type deleteQueueCmd struct {
	*cobra.Command
	opts *deleteQueueOpts
}

func newDeleteQueueCmd(cfg *client.Config) *deleteQueueCmd {
	opts := &deleteQueueOpts{}
	c := &deleteQueueCmd{
		Command: &cobra.Command{
			Use:     "queue NAME",
			Args:    cobra.ExactArgs(1),
			Aliases: []string{"q"},
		},
		opts: opts,
	}

	// cmd := c.Cmd()
	// flags := cmd.Flags()

	return c
}

func (c *deleteQueueCmd) Cmd() *cobra.Command { return c.Command }
func (c *deleteQueueCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	return runDeleteQueue(ctx, cfg, c.opts, cmd, args)
}

func runDeleteQueue(ctx context.Context, cfg *client.Config, opts *deleteQueueOpts, cmd *cobra.Command, args []string) error {
	name := args[0]
	cl := clientFromContext(ctx)
	return cl.DeleteQueue(ctx, name)
}
