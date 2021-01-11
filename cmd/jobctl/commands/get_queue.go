package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
	// apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
)

type getQueueOpts struct {
	selectorRaw string
}

type getQueueCmd struct {
	*cobra.Command
	opts *getQueueOpts
}

func newGetQueueCmd(cfg *client.Config) *getQueueCmd {
	opts := &getQueueOpts{}
	c := &getQueueCmd{
		opts: opts,
		Command: &cobra.Command{
			Use:     "queue",
			Aliases: []string{"q"},
			Args:    cobra.ExactArgs(1),
		},
	}

	return c
}

func (c *getQueueCmd) Cmd() *cobra.Command { return c.Command }
func (c *getQueueCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	cl := clientFromContext(ctx)
	q, err := cl.GetQueue(ctx, args[0])
	if err != nil {
		return err
	}
	fmt.Println(q)
	return nil
}
