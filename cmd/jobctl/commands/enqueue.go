package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
	"github.com/jeffrom/job-manager/pkg/config"
)

type enqueueOpts struct {
	// jobclient.EnqueueOpts
}

type enqueueCmd struct {
	*cobra.Command
	opts *enqueueOpts
}

func newEnqueueCmd(cfg *jobclient.Config) *enqueueCmd {
	opts := &enqueueOpts{}
	c := &enqueueCmd{
		Command: &cobra.Command{
			Use:     "enqueue",
			Args:    cobra.MinimumNArgs(1),
			Aliases: []string{"enq"},
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				ctx := cmd.Context()
				if ctx == nil {
					ctx = context.Background()
				}
				if len(args) == 0 {
					return handleCompletion(completeQueueList(ctx, toComplete))
				}
				return nil, cobra.ShellCompDirectiveNoFileComp
			},
		},
		opts: opts,
	}

	return c
}

func (c *enqueueCmd) Cmd() *cobra.Command { return c.Command }
func (c *enqueueCmd) Execute(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
	client := clientFromContext(ctx)
	iargs := argsToInterface(args[1:])
	id, err := client.EnqueueJob(ctx, args[0], iargs...)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", id)
	return nil
}

func handleCompletion(comps []string, err error) ([]string, cobra.ShellCompDirective) {
	f, _ := os.Create("hi.txt")
	defer f.Close()
	fmt.Fprintf(f, "comps: %+v, err: %v\n", comps, err)
	if err != nil {
		// TODO maybe a better way to print the error
		fmt.Fprintf(os.Stderr, "completion error: %v\n", err)
		return comps, cobra.ShellCompDirectiveError
	}
	return comps, cobra.ShellCompDirectiveNoFileComp
}

func completeQueueList(ctx context.Context, toComplete string) ([]string, error) {
	icfg, err := config.MergeEnvFlags(&jobclient.Config{}, &jobclient.ConfigDefaults)
	if err != nil {
		return nil, err
	}
	cfg := icfg.(*jobclient.Config)

	client := jobclient.New(cfg.Addr, jobclient.WithConfig(cfg))
	queues, err := client.ListQueues(ctx, jobclient.ListQueuesOpts{})
	if err != nil {
		return nil, err
	}

	names := make([]string, len(queues.Queues))
	for i, q := range queues.Queues {
		names[i] = q.ID
	}
	return names, nil
}

func argsToInterface(args []string) []interface{} {
	ifaces := make([]interface{}, len(args))
	for i, arg := range args {
		ifaces[i] = arg
	}
	return ifaces
}
