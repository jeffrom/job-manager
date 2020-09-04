package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/jobclient"
	"github.com/jeffrom/job-manager/pkg/label"
	"github.com/jeffrom/job-manager/pkg/schema"
)

type saveQueueOpts struct {
	jobclient.SaveQueueOpts

	LabelFlags []string
	SchemaPath string
}

type saveQueueCmd struct {
	*cobra.Command
	opts *saveQueueOpts
}

func newSaveQueueCmd(cfg *jobclient.Config) *saveQueueCmd {
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
	flags.StringArrayVarP(&opts.LabelFlags, "label", "l", nil, "set label `name=value`")
	flags.StringVarP(&opts.SchemaPath, "schema", "S", "", "path to json schema")

	return c
}

func (c *saveQueueCmd) Cmd() *cobra.Command { return c.Command }
func (c *saveQueueCmd) Execute(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
	return runSaveQueue(ctx, cfg, c.opts, cmd, args)
}

func runSaveQueue(ctx context.Context, cfg *jobclient.Config, opts *saveQueueOpts, cmd *cobra.Command, args []string) error {
	labels, err := label.ParseStringArray(opts.LabelFlags)
	if err != nil {
		return err
	}
	// TODO reading the schema from stdin could be cool too
	var scm *schema.Schema
	if p := opts.SchemaPath; p != "" {
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		scm, err = schema.ParseBytes(b)
		if err != nil {
			return err
		}
	}
	var argSchema []byte
	var dataSchema []byte
	var resultSchema []byte
	if scm != nil {
		argSchema, err = json.Marshal(scm.Args)
		if err != nil {
			return err
		}
		dataSchema, err = json.Marshal(scm.Data)
		if err != nil {
			return err
		}
		resultSchema, err = json.Marshal(scm.Result)
		if err != nil {
			return err
		}
	}

	client := clientFromContext(ctx)
	q, err := client.SaveQueue(ctx, args[0], jobclient.SaveQueueOpts{
		Concurrency:  opts.Concurrency,
		MaxRetries:   opts.MaxRetries,
		JobDuration:  opts.JobDuration,
		Labels:       labels,
		ArgSchema:    argSchema,
		DataSchema:   dataSchema,
		ResultSchema: resultSchema,
	})
	if err != nil {
		return err
	}
	fmt.Printf("<- %+v\n", q)
	return nil
}
