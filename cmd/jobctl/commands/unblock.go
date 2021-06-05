package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
)

type unblockCmdOpts struct {
}

type unblockCmd struct {
	*cobra.Command
	opts *unblockCmdOpts
}

func newUnblockCmd(cfg *client.Config) *unblockCmd {
	opts := &unblockCmdOpts{}
	c := &unblockCmd{
		Command: &cobra.Command{
			Use:               "unblock",
			Args:              cobra.ExactArgs(1),
			ValidArgsFunction: validQueueList(1),
		},
		opts: opts,
	}

	// cmd := c.Cmd()
	// flags := cmd.Flags()

	return c
}

func (c *unblockCmd) Cmd() *cobra.Command { return c.Command }
func (c *unblockCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	client := clientFromContext(ctx)
	if err := client.UnblockQueue(ctx, args[0]); err != nil {
		return err
	}
	fmt.Println("ok")
	return nil
}
