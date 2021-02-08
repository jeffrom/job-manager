package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
)

type statsCmd struct {
	*cobra.Command
}

func newStatsCmd(cfg *client.Config) *statsCmd {
	c := &statsCmd{
		Command: &cobra.Command{
			Use: "stats",
		},
	}

	return c
}

func (c *statsCmd) Cmd() *cobra.Command { return c.Command }
func (c *statsCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	client := clientFromContext(ctx)
	stats, err := client.Stats(ctx)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
