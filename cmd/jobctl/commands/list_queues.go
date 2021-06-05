package commands

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/label"
	"github.com/jeffrom/job-manager/mjob/resource"
	// apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
)

type listQueuesOpts struct {
	selectorRaw string
	limit       int64
	lastID      string
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
	flags.Int64VarP(&opts.limit, "limit", "L", 20, "per-page limit")
	flags.StringVarP(&opts.lastID, "last-id", "l", "", "last id (from previous page)")

	return c
}

func (c *listQueuesCmd) Cmd() *cobra.Command { return c.Command }
func (c *listQueuesCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	cl := clientFromContext(ctx)
	qs, err := cl.ListQueues(ctx, client.ListQueuesOpts{
		Names:     args,
		Selectors: label.SplitSelectors(c.opts.selectorRaw),
		Page: &resource.Pagination{
			Limit:  c.opts.limit,
			LastID: c.opts.lastID,
		},
	})
	if err != nil {
		return err
	}
	padding := 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "NAME\tCREATED\tVERSION\tPAUSED\tBLOCKED\tRETRIES\tDURATION\tBACKOFF\tUNIQUE\tLABELS\t")
	for _, q := range qs.Queues {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%d\t%s\t%s\t%s\t%s\n",
			q.Name,
			q.CreatedAt.Local().Format(time.Stamp),
			q.Version.String(),
			yesno(q.Paused),
			yesno(q.Blocked),
			q.Retries,
			cleanupDuration(q.Duration.String()),
			queueBackoff(q),
			yesno(q.Unique),
			q.Labels.String())
	}
	// fmt.Fprintln(w)
	return w.Flush()
}

func yesno(b bool) string {
	if b {
		return "y"
	}
	return ""
}

func queueBackoff(q *resource.Queue) string {
	if q.BackoffInitial == 0 || q.BackoffFactor == 0 {
		return ""
	}

	return fmt.Sprintf("%s<>%s x %.2f",
		cleanupDuration(q.BackoffInitial.String()),
		cleanupDuration(q.BackoffMax.String()),
		q.BackoffFactor)
}

var durationCleanupRe = regexp.MustCompile(`([a-z])0[a-z]$`)

func cleanupDuration(d string) string {
	return durationCleanupRe.ReplaceAllString(d, "$1")
}
