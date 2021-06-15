package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/pkg/config"
)

func handleCompletion(comps []string, err error) ([]string, cobra.ShellCompDirective) {
	// f, _ := os.Create("hi.txt")
	// defer f.Close()
	// fmt.Fprintf(f, "comps: %+v, err: %v\n", comps, err)
	// TODO better way to print errors, debug info, need to dump to txt file
	if err != nil {
		// fmt.Fprintf(os.Stderr, "completion error: %v\n", err)
		return comps, cobra.ShellCompDirectiveError
	}
	return comps, cobra.ShellCompDirectiveNoFileComp
}

func validQueueList(n int) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		if n < 0 || len(args) < n {
			return handleCompletion(completeQueueList(ctx, toComplete))
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
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
