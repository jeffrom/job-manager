package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/label"
	"github.com/jeffrom/job-manager/mjob/resource"
	"github.com/jeffrom/job-manager/pkg/config"
)

type enqueueOpts struct {
	// client.EnqueueOpts
	// data string
	claims []string
}

type enqueueCmd struct {
	*cobra.Command
	opts *enqueueOpts
}

func newEnqueueCmd(cfg *client.Config) *enqueueCmd {
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

	flags := c.Command.Flags()
	flags.StringArrayVarP(&opts.claims, "claim", "c", nil, "enqueue with claims")

	return c
}

func (c *enqueueCmd) Cmd() *cobra.Command { return c.Command }
func (c *enqueueCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	var claims label.Claims
	if len(c.opts.claims) > 0 {
		var err error
		claims, err = label.ParseClaims(c.opts.claims)
		if err != nil {
			return err
		}
	}
	cl := clientFromContext(ctx)
	iargs, err := resource.ParseCLIArgs(args[1:])
	if err != nil {
		return err
	}
	opts := client.EnqueueOpts{
		Claims: claims,
	}
	id, err := cl.EnqueueJobOpts(ctx, args[0], opts, iargs...)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", id)
	return nil
}

func handleCompletion(comps []string, err error) ([]string, cobra.ShellCompDirective) {
	// f, _ := os.Create("hi.txt")
	// defer f.Close()
	// fmt.Fprintf(f, "comps: %+v, err: %v\n", comps, err)
	if err != nil {
		// TODO better way to print errors, debug info, need to dump to txt file
		// fmt.Fprintf(os.Stderr, "completion error: %v\n", err)
		return comps, cobra.ShellCompDirectiveError
	}
	return comps, cobra.ShellCompDirectiveNoFileComp
}

func completeQueueList(ctx context.Context, toComplete string) ([]string, error) {
	icfg, err := config.MergeEnvFlags(&client.Config{}, &client.ConfigDefaults)
	if err != nil {
		return nil, err
	}
	cfg := icfg.(*client.Config)

	c := client.New(cfg.Host, client.WithConfig(cfg))
	queues, err := c.ListQueues(ctx, client.ListQueuesOpts{})
	if err != nil {
		return nil, err
	}

	names := make([]string, len(queues.Queues))
	for i, q := range queues.Queues {
		names[i] = q.Name
	}
	return names, nil
}
