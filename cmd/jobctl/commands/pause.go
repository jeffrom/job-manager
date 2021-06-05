package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
)

type pauseCmdOpts struct {
}

type pauseCmd struct {
	*cobra.Command
	opts *pauseCmdOpts
}

func newPauseCmd(cfg *client.Config) *pauseCmd {
	opts := &pauseCmdOpts{}
	c := &pauseCmd{
		Command: &cobra.Command{
			Use:               "pause",
			Args:              cobra.ExactArgs(1),
			ValidArgsFunction: validQueueList(1),
		},
		opts: opts,
	}

	// cmd := c.Cmd()
	// flags := cmd.Flags()

	return c
}

func (c *pauseCmd) Cmd() *cobra.Command { return c.Command }
func (c *pauseCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	client := clientFromContext(ctx)
	if err := client.PauseQueue(ctx, args[0]); err != nil {
		return err
	}
	fmt.Println("ok")
	return nil
}
