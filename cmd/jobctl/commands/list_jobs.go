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
	"github.com/jeffrom/job-manager/mjob/resource"
	// apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
)

type listJobsOpts struct {
	selectorRaw string
	statuses    []string
	limit       int64
	lastID      string
}

type listJobsCmd struct {
	*cobra.Command
	opts *listJobsOpts
}

func newListJobsCmd(cfg *client.Config) *listJobsCmd {
	opts := &listJobsOpts{}
	c := &listJobsCmd{
		opts: opts,
		Command: &cobra.Command{
			Use:     "jobs",
			Aliases: []string{"job", "j"},
			Args:    cobra.ArbitraryArgs,
		},
	}

	cmd := c.Cmd()
	flags := cmd.Flags()
	flags.StringVarP(&opts.selectorRaw, "selector", "s", "", "filter by selector")
	flags.StringArrayVarP(&opts.statuses, "status", "S", nil, "filter by status")
	flags.Int64VarP(&opts.limit, "limit", "L", 20, "per-page limit")
	flags.StringVarP(&opts.lastID, "last-id", "l", "", "last id (from previous page)")

	cmd.RegisterFlagCompletionFunc("status", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{
			"queued", "running", "complete", "failed", "dead", "invalid", "cancelled",
		}, cobra.ShellCompDirectiveDefault
	})

	return c
}

func (c *listJobsCmd) Cmd() *cobra.Command { return c.Command }

func (c *listJobsCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	cl := clientFromContext(ctx)
	statuses, err := statusFromStrings(c.opts.statuses)
	if err != nil {
		return err
	}
	jobs, err := cl.ListJobs(ctx, client.ListJobsOpts{
		Queues:    args,
		Selectors: label.SplitSelectors(c.opts.selectorRaw),
		Statuses:  statuses,
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
	fmt.Fprintln(w, "ID\tNAME\tENQUEUED\tVERSION/Q\tATTEMPTS\tSTATUS\t")
	for _, jb := range jobs.Jobs {
		// arg, _ := json.Marshal(jb.Args)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\t\n",
			jb.ID,
			jb.Name,
			jb.EnqueuedAt.Local().Format(time.Stamp),
			fmt.Sprintf("%s/%s", jb.Version.String(), jb.QueueVersion.String()),
			jb.Attempt,
			jb.Status.String(),
			// string(arg),
		)
	}
	return w.Flush()
}

func statusFromStrings(statuses []string) ([]resource.Status, error) {
	res := make([]resource.Status, len(statuses))
	for i, st := range statuses {
		res[i] = *resource.StatusFromString(st)
	}
	return res, nil
}
