// Package commands contains jobctl's cobra commands.
package commands

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
	"github.com/jeffrom/job-manager/pkg/config"
)

func ExecuteArgs(args []string) error {
	cfg := &jobclient.Config{}
	cmd := newRootCmd(cfg)
	cmd.SetArgs(args)
	ctx := context.Background()
	if err := cmd.ExecuteContext(ctx); err != nil {
		// fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}

func Execute() {
	args := os.Args[1:]
	if err := ExecuteArgs(args); err != nil {
		os.Exit(1)
	}
}

type wrappedCobraRun = func(cfg *jobclient.Config, cmd *cobra.Command, args []string) error

func wrapCmd(cfgFlags *jobclient.Config, fn wrappedCobraRun) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		icfg, err := config.MergeEnvFlags(cfgFlags, &jobclient.ConfigDefaults)
		if err != nil {
			return err
		}
		return fn(icfg.(*jobclient.Config), cmd, args)
	}
}
