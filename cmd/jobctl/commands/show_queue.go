package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
	// apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
)

type showQueueOpts struct {
	// selectorRaw string
}

type showQueueCmd struct {
	*cobra.Command
	opts *showQueueOpts
}

func newShowQueueCmd(cfg *client.Config) *showQueueCmd {
	opts := &showQueueOpts{}
	c := &showQueueCmd{
		opts: opts,
		Command: &cobra.Command{
			Use:     "queue",
			Aliases: []string{"q"},
			Args:    cobra.ExactArgs(1),
		},
	}

	return c
}

func (c *showQueueCmd) Cmd() *cobra.Command { return c.Command }
func (c *showQueueCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	cl := clientFromContext(ctx)
	q, err := cl.GetQueue(ctx, args[0])
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
