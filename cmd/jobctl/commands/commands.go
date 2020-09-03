// Package commands contains jobctl's cobra commands.
package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
	"github.com/jeffrom/job-manager/pkg/config"
)

type Command interface {
	Cmd() *cobra.Command
	Execute(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error
}

type contextKey string

var ctxClientKey contextKey = "client"

func ExecuteArgs(args []string) error {
	cfg := &jobclient.Config{}
	cmd := newRootCmd(cfg)
	cmd.SetArgs(args)
	ctx := context.Background()
	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
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

type wrappedCobraRun = func(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error

func wrapCmdRun(cfgFlags *jobclient.Config, fn wrappedCobraRun) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		icfg, err := config.MergeEnvFlags(cfgFlags, &jobclient.ConfigDefaults)
		if err != nil {
			return err
		}
		cfg := icfg.(*jobclient.Config)

		ctx := cmd.Context()
		c := jobclient.New(cfg.Addr, jobclient.WithConfig(cfg))
		ctx = context.WithValue(ctx, ctxClientKey, c)
		return fn(ctx, cfg, cmd, args)
	}
}

func clientFromContext(ctx context.Context) *jobclient.Client {
	return ctx.Value(ctxClientKey).(*jobclient.Client)
}

func addCommands(cfg *jobclient.Config, parent Command, children ...Command) {
	p := parent.Cmd()
	for _, child := range children {
		childCmd := child.Cmd()
		childCmd.RunE = wrapCmdRun(cfg, child.Execute)
		p.AddCommand(childCmd)
	}
}
