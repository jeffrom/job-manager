package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/label"
	"github.com/jeffrom/job-manager/mjob/schema"
)

type saveQueueOpts struct {
	client.SaveQueueOpts

	LabelFlags []string
	SchemaPath string

	BackoffInitial time.Duration
	BackoffMax     time.Duration
	BackoffFactor  float32
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
	flags.IntVarP(&opts.MaxRetries, "retries", "r", 0, "max retries")
	flags.DurationVarP(&opts.JobDuration, "duration", "d", 0, "job max duration")
	flags.DurationVar(&opts.CheckinDuration, "checkin-duration", 0, "job checkin duration")
	flags.DurationVar(&opts.ClaimDuration, "claim-duration", 0, "job claim duration")
	flags.StringArrayVarP(&opts.LabelFlags, "label", "l", nil, "set label `name=value`")
	flags.StringVarP(&opts.SchemaPath, "schema", "S", "", "path to json schema")
	flags.BoolVarP(&opts.Unique, "unique", "U", false, "run one unique arg list concurrently")
	flags.DurationVar(&opts.BackoffInitial, "backoff-initial", 0, "initial backoff duration")
	flags.DurationVar(&opts.BackoffMax, "backoff-max", 0, "max backoff duration")
	flags.Float32Var(&opts.BackoffFactor, "backoff-factor", 1.0, "backoff factor")

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

	var boInitial time.Duration
	var boMax time.Duration
	var boFactor float32
	if opts.BackoffFactor > 0 {
		boFactor = opts.BackoffFactor
	}
	if opts.BackoffInitial > 0 {
		boInitial = opts.BackoffInitial
	}
	if opts.BackoffMax > 0 {
		boMax = opts.BackoffMax
	}

	q, err := cl.SaveQueue(ctx, id, client.SaveQueueOpts{
		MaxRetries:      opts.MaxRetries,
		JobDuration:     opts.JobDuration,
		CheckinDuration: opts.CheckinDuration,
		ClaimDuration:   opts.ClaimDuration,
		Labels:          labels,
		Schema:          scmb,
		Unique:          opts.Unique,
		Version:         v,
		BackoffInitial:  boInitial,
		BackoffMax:      boMax,
		BackoffFactor:   boFactor,
	})
	if err != nil {
		return err
	}
	fmt.Printf("<- %+v\n", q)
	return nil
}
