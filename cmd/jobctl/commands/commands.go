// Package commands contains jobctl's cobra commands.
package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/pkg/config"
)

type Command interface {
	Cmd() *cobra.Command
	Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error
}

type contextKey string

var ctxClientKey contextKey = "client"

func ExecuteArgs(args []string) error {
	cfg := &client.Config{}
	if host := os.Getenv("HOST"); host != "" {
		cfg.Host = host
	}
	cmd := newRootCmd(cfg)
	cmd.SetArgs(args)
	ctx := context.Background()
	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
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

func wrapCobraCommand(cfg *client.Config, c Command) *cobra.Command {
	cmd := c.Cmd()
	cmd.RunE = wrapCmdRun(cfg, c.Execute)
	return cmd
}

type wrappedCobraRun = func(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error

func wrapCmdRun(cfgFlags *client.Config, fn wrappedCobraRun) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		icfg, err := config.MergeEnvFlags(cfgFlags, &client.ConfigDefaults)
		if err != nil {
			return err
		}
		cfg := icfg.(*client.Config)

		ctx := cmd.Context()
		c := client.New(cfg.Host,
			client.WithConfig(cfg),
			client.WithHTTPClient(httpClient),
		)
		ctx = context.WithValue(ctx, ctxClientKey, c)
		return fn(ctx, cfg, cmd, args)
	}
}

func clientFromContext(ctx context.Context) *client.Client {
	return ctx.Value(ctxClientKey).(*client.Client)
}

func addCommands(cfg *client.Config, parent Command, children ...Command) {
	p := parent.Cmd()
	for _, child := range children {
		childCmd := child.Cmd()
		childCmd.RunE = wrapCmdRun(cfg, child.Execute)
		p.AddCommand(childCmd)
	}
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}
