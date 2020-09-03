package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
)

func newQueueListCmd(cfg *jobclient.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		RunE:    wrapCmd(cfg, queueListCmd),
	}

	return cmd
}

func queueListCmd(cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
	fmt.Printf("%+v\n", cfg)
	return nil
}
