package commands

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/pkg/label"
	"github.com/jeffrom/job-manager/pkg/schema"
)

type saveQueueOpts struct {
	client.SaveQueueOpts

	LabelFlags []string
	SchemaPath string
}

type saveQueueCmd struct {
	*cobra.Command
	opts *saveQueueOpts
}

func newSaveQueueCmd(cfg *client.Config) *saveQueueCmd {
	opts := &saveQueueOpts{}
	c := &saveQueueCmd{
		Command: &cobra.Command{
			Use:     "queue",
			Args:    cobra.ExactArgs(1),
			Aliases: []string{"q"},
		},
		opts: opts,
	}

	cmd := c.Cmd()
	flags := cmd.Flags()
	flags.IntVarP(&opts.Concurrency, "concurrency", "c", 0, "job concurrency")
	flags.IntVarP(&opts.MaxRetries, "retries", "r", 0, "max retries")
	flags.DurationVarP(&opts.JobDuration, "duration", "d", 0, "job max duration")
	flags.DurationVar(&opts.CheckinDuration, "checkin-duration", 0, "job checkin duration")
	flags.DurationVar(&opts.ClaimDuration, "claim-duration", 0, "job claim duration")
	flags.StringArrayVarP(&opts.LabelFlags, "label", "l", nil, "set label `name=value`")
	flags.StringVarP(&opts.SchemaPath, "schema", "S", "", "path to json schema")
	flags.BoolVarP(&opts.Unique, "unique", "U", false, "run one unique arg list concurrently")

	return c
}

func (c *saveQueueCmd) Cmd() *cobra.Command { return c.Command }
func (c *saveQueueCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	return runSaveQueue(ctx, cfg, c.opts, cmd, args)
}

func runSaveQueue(ctx context.Context, cfg *client.Config, opts *saveQueueOpts, cmd *cobra.Command, args []string) error {
	labels, err := label.ParseStringArray(opts.LabelFlags)
	if err != nil {
		return err
	}
	// TODO reading the schema, or whole queue cfg from stdin could be cool too
	var scmb []byte
	if p := opts.SchemaPath; p != "" {
		var err error
		scmb, err = ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		_, err = schema.Parse(scmb)
		if err != nil {
			return err
		}
	}

	id := args[0]
	cl := clientFromContext(ctx)
	prev, err := cl.GetQueue(ctx, id)
	if err != nil && !client.IsNotFound(err) {
		return err
	}
	v := ""
	if prev != nil {
		v = prev.Version.Strict()
	}

	q, err := cl.SaveQueue(ctx, id, client.SaveQueueOpts{
		Concurrency:     opts.Concurrency,
		MaxRetries:      opts.MaxRetries,
		JobDuration:     opts.JobDuration,
		CheckinDuration: opts.CheckinDuration,
		ClaimDuration:   opts.ClaimDuration,
		Labels:          labels,
		Schema:          scmb,
		Unique:          opts.Unique,
		Version:         v,
	})
	if err != nil {
		return err
	}
	fmt.Printf("<- %+v\n", q)
	return nil
}
