package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
)

type rootCmd struct {
	*cobra.Command
}

func (c *rootCmd) Cmd() *cobra.Command { return c.Command }

func newRootCmd(cfg *client.Config) *rootCmd {
	c := &rootCmd{
		Command: &cobra.Command{
			Use:           "jobctl",
			SilenceErrors: true, // we are printing errors ourselves
			SilenceUsage:  true,
			Args:          cobra.NoArgs,
			RunE:          wrapCmdRun(cfg, usageCmd),
		},
	}
	cmd := c.Cmd()

	flags := cmd.PersistentFlags()
	flags.StringVarP(&cfg.Host, "host", "H", "", "set host:port (env: $HOST)")

	cmd.AddCommand(
		wrapCobraCommand(cfg, newListCmd(cfg)),
		wrapCobraCommand(cfg, newSaveCmd(cfg)),
		wrapCobraCommand(cfg, newEnqueueCmd(cfg)),
		wrapCobraCommand(cfg, newAckCmd(cfg)),
		wrapCobraCommand(cfg, newConsumerCmd(cfg)),
		wrapCobraCommand(cfg, newCompletionCmd(cfg)),
	)
	return c
}

func usageCmd(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
