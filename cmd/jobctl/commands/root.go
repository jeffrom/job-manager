package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
)

type rootCmd struct {
	*cobra.Command
}

func (c *rootCmd) Cmd() *cobra.Command { return c.Command }

func newRootCmd(cfg *jobclient.Config) *rootCmd {
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
	flags.StringVarP(&cfg.Addr, "host", "H", "", "set host:port (env: $HOST)")

	cmd.AddCommand(
		newListCmd(cfg).Cmd(),
		newSaveCmd(cfg).Cmd(),
	)
	return c
}

func usageCmd(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
