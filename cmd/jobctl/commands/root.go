package commands

import (
	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
)

func newRootCmd(cfg *jobclient.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use: "jobctl",
	}

	flags := cmd.PersistentFlags()
	flags.StringVarP(&cfg.Addr, "host", "H", "", "set host:port (env: $HOST)")

	cmd.AddCommand(
		newQueueCmd(cfg),
	)
	return cmd
}
