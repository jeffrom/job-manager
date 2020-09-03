package commands

import (
	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
)

func newQueueCmd(cfg *jobclient.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "queue",
		Aliases: []string{"q"},
	}

	return cmd
}
