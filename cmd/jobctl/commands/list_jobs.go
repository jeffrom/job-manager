package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/label"
	// apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
)

type listJobsOpts struct {
	selectorRaw string
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

	return c
}

func (c *listJobsCmd) Cmd() *cobra.Command { return c.Command }
func (c *listJobsCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	cl := clientFromContext(ctx)
	jobs, err := cl.ListJobs(ctx, client.ListJobsOpts{
		Queues:    args,
		Selectors: label.SplitSelectors(c.opts.selectorRaw),
	})
	if err != nil {
		return err
	}
	padding := 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "NAME\tENQUEUED\tVERSION/Q\tATTEMPTS\tSTATUS\tARGS\t")
	for _, jb := range jobs.Jobs {
		arg, _ := json.Marshal(jb.Args)
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\t\n",
			jb.Name,
			jb.EnqueuedAt.Local().Format(time.Stamp),
			fmt.Sprintf("%s/%s", jb.Version.String(), jb.QueueVersion.String()),
			jb.Attempt,
			jb.Status.String(),
			string(arg),
		)
	}
	return w.Flush()
}
