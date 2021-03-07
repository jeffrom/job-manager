package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/resource"
)

type ackCmdOpts struct {
	cancel  bool
	fail    bool
	invalid bool
}

type ackCmd struct {
	*cobra.Command
	opts *ackCmdOpts
}

func newAckCmd(cfg *client.Config) *ackCmd {
	opts := &ackCmdOpts{}
	c := &ackCmd{
		Command: &cobra.Command{
			Use: "ack",
		},
		opts: opts,
	}

	cmd := c.Cmd()
	flags := cmd.Flags()
	flags.BoolVarP(&opts.cancel, "cancel", "c", false, "cancel job")
	flags.BoolVarP(&opts.fail, "fail", "f", false, "fail job")
	flags.BoolVarP(&opts.invalid, "invalid", "i", false, "mark job invalid")

	return c
}

func (c *ackCmd) Cmd() *cobra.Command { return c.Command }
func (c *ackCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	client := clientFromContext(ctx)
	status := resource.StatusComplete
	if c.opts.cancel {
		status = resource.StatusCancelled
	} else if c.opts.fail {
		status = resource.StatusFailed
	} else if c.opts.invalid {
		status = resource.StatusInvalid
	}
	return client.AckJob(ctx, args[0], status)
}
