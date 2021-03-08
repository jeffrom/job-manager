package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/apply"
	"github.com/jeffrom/job-manager/mjob/client"
)

type applyOpts struct {
	// path string
}

type applyCmd struct {
	*cobra.Command
	opts *applyOpts
}

func newApplyCmd(cfg *client.Config) *applyCmd {
	opts := &applyOpts{}
	c := &applyCmd{
		Command: &cobra.Command{
			Use:  "apply",
			Args: cobra.ArbitraryArgs,
		},
		opts: opts,
	}

	// flags := c.Cmd().Flags()
	// flags.StringVarP(&opts.path, "filename", "f", "", "path to a queue resource yaml file")
	return c
}

func (c *applyCmd) Cmd() *cobra.Command { return c.Command }
func (c *applyCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	cl := clientFromContext(ctx)
	for _, arg := range args {
		if err := apply.Path(ctx, cl, arg); err != nil {
			return err
		}
	}
	return nil
}
