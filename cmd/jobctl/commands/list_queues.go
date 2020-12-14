package commands

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/label"
	// apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
)

type listQueuesOpts struct {
	selectorRaw string
}

type listQueuesCmd struct {
	*cobra.Command
	opts *listQueuesOpts
}

func newListQueuesCmd(cfg *client.Config) *listQueuesCmd {
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
func (c *listQueuesCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	cl := clientFromContext(ctx)
	qs, err := cl.ListQueues(ctx, client.ListQueuesOpts{
		Names:     args,
		Selectors: label.SplitSelectors(c.opts.selectorRaw),
	})
	if err != nil {
		return err
	}
	padding := 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "NAME\tCREATED\tVERSION\tLABELS\t")
	for _, q := range qs.Queues {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", q.Name, q.CreatedAt.Format(time.Stamp), q.Version.String(), q.Labels.String())
	}
	// fmt.Fprintln(w)
	return w.Flush()
}
