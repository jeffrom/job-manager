package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
)

type unpauseCmdOpts struct {
}

type unpauseCmd struct {
	*cobra.Command
	opts *unpauseCmdOpts
}

func newUnpauseCmd(cfg *client.Config) *unpauseCmd {
	opts := &unpauseCmdOpts{}
	c := &unpauseCmd{
		Command: &cobra.Command{
			Use:               "unpause",
			Args:              cobra.ExactArgs(1),
			ValidArgsFunction: validQueueList(1),
		},
		opts: opts,
	}

	// cmd := c.Cmd()
	// flags := cmd.Flags()

	return c
}

func (c *unpauseCmd) Cmd() *cobra.Command { return c.Command }
func (c *unpauseCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	client := clientFromContext(ctx)
	if err := client.UnpauseQueue(ctx, args[0]); err != nil {
		return err
	}
	fmt.Println("ok")
	return nil
}
