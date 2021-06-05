package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
)

type showJobOpts struct {
	includes []string
}

type showJobCmd struct {
	*cobra.Command
	opts *showJobOpts
}

func newShowJobCmd(cfg *client.Config) *showJobCmd {
	opts := &showJobOpts{}
	c := &showJobCmd{
		opts: opts,
		Command: &cobra.Command{
			Use:     "job",
			Aliases: []string{"j"},
			Args:    cobra.ExactArgs(1),
		},
	}

	cmd := c.Cmd()
	flags := cmd.Flags()
	flags.StringArrayVarP(&opts.includes, "include", "i", allowedJobIncludes, fmt.Sprintf("include additional data (%q)", allowedJobIncludes))
	die(cmd.RegisterFlagCompletionFunc("include", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"checkin", "result"}, cobra.ShellCompDirectiveNoFileComp
	}))

	return c
}

func (c *showJobCmd) Cmd() *cobra.Command { return c.Command }
func (c *showJobCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	cl := clientFromContext(ctx)
	job, err := cl.GetJob(ctx, args[0], &client.GetJobOpts{
		Includes: c.opts.includes,
	})
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(job, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
