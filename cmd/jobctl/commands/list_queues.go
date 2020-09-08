package commands

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
	"github.com/jeffrom/job-manager/pkg/label"
	// apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
)

type listQueuesOpts struct {
	selectorRaw string
}

type listQueuesCmd struct {
	*cobra.Command
	opts *listQueuesOpts
}

func newListQueuesCmd(cfg *jobclient.Config) *listQueuesCmd {
	opts := &listQueuesOpts{}
	c := &listQueuesCmd{
		opts: opts,
		Command: &cobra.Command{
			Use:     "queues",
			Aliases: []string{"queue", "q"},
			Args:    cobra.ArbitraryArgs,
		},
	}

	cmd := c.Cmd()
	flags := cmd.Flags()
	flags.StringVarP(&opts.selectorRaw, "selector", "s", "", "filter by selector")

	return c
}

func (c *listQueuesCmd) Cmd() *cobra.Command { return c.Command }
func (c *listQueuesCmd) Execute(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
	client := clientFromContext(ctx)
	queues, err := client.ListQueues(ctx, jobclient.ListQueuesOpts{
		Names:     args,
		Selectors: label.SplitSelectors(c.opts.selectorRaw),
	})
	if err != nil {
		return err
	}
	padding := 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "NAME\tCREATED\tVERSION\t")
	for _, q := range queues.Queues {
		fmt.Fprintf(w, "%s\t%s\t%d\n", q.ID, q.CreatedAt.Format(time.Stamp), q.Version)
	}
	// fmt.Fprintln(w)
	return w.Flush()
}

// func runListQueues(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
// 	client := clientFromContext(ctx)
// 	queues, err := client.ListQueues(ctx, jobclient.ListQueuesOpts{})
// 	if err != nil {
// 		return err
// 	}
// 	padding := 3
// 	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
// 	fmt.Fprintln(w, "NAME\tCREATED\t")
// 	for _, q := range queues.Queues {
// 		fmt.Fprintf(w, "%s\t%s\n", q.Id, q.CreatedAt.AsTime().Format(time.Stamp))
// 	}
// 	// fmt.Fprintln(w)
// 	return w.Flush()
// }
