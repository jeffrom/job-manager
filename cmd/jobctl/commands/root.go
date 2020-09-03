package commands

import (
	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
)

func newRootCmd(cfg *jobclient.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "jobctl",
		SilenceUsage: true,
		RunE:         wrapCmd(cfg, usageCmd),
	}

	flags := cmd.PersistentFlags()
	flags.StringVarP(&cfg.Addr, "host", "H", "", "set host:port (env: $HOST)")

	cmd.AddCommand(
		newQueueCmd(cfg),
	)
	return cmd
}

func usageCmd(cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
