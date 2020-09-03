package commands

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

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
	client := clientFromContext(ctx)
	queues, err := client.ListQueues(ctx, jobclient.ListQueuesOpts{})
	if err != nil {
		return err
	}
	padding := 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "NAME\tCREATED\t")
	for _, q := range queues.Queues {
		fmt.Fprintf(w, "%s\t%s\n", q.Id, q.CreatedAt.AsTime().Format(time.Stamp))
	}
	// fmt.Fprintln(w)
	return w.Flush()
}
