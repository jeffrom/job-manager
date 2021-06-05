package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
)

type blockCmdOpts struct {
}

type blockCmd struct {
	*cobra.Command
	opts *blockCmdOpts
}

func newBlockCmd(cfg *client.Config) *blockCmd {
	opts := &blockCmdOpts{}
	c := &blockCmd{
		Command: &cobra.Command{
			Use:               "block",
			Args:              cobra.ExactArgs(1),
			ValidArgsFunction: validQueueList(1),
		},
		opts: opts,
	}

	// cmd := c.Cmd()
	// flags := cmd.Flags()

	return c
}

func (c *blockCmd) Cmd() *cobra.Command { return c.Command }
func (c *blockCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	client := clientFromContext(ctx)
	if err := client.BlockQueue(ctx, args[0]); err != nil {
		return err
	}
	fmt.Println("ok")
	return nil
}
