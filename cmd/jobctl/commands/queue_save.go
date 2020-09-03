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

type queueSaveOpts struct {
	jobclient.SaveQueueOpts

	LabelFlags []string
	SchemaPath string
}

type queueSaveCmd struct {
	*cobra.Command
	opts *queueSaveOpts
}

func newQueueSaveCmd(cfg *jobclient.Config) *queueSaveCmd {
	opts := &queueSaveOpts{}
	c := &queueSaveCmd{
		Command: &cobra.Command{
			Use:  "save",
			Args: cobra.ExactArgs(1),
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

func (c *queueSaveCmd) Cmd() *cobra.Command { return c.Command }
func (c *queueSaveCmd) Execute(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
	labels, err := label.ParseStringArray(c.opts.LabelFlags)
	if err != nil {
		return err
	}
	var scm *schema.Schema
	if p := c.opts.SchemaPath; p != "" {
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
		Concurrency:  c.opts.Concurrency,
		MaxRetries:   c.opts.MaxRetries,
		JobDuration:  c.opts.JobDuration,
		Labels:       labels,
		ArgSchema:    argSchema,
		DataSchema:   dataSchema,
		ResultSchema: resultSchema,
	})
	if err != nil {
		return err
	}
	fmt.Printf("res: %+v\n", q)
	return nil
}

func queueSaveRun(ctx context.Context, cfg *jobclient.Config, cmd *cobra.Command, args []string) error {
	c := clientFromContext(ctx)
	q, err := c.SaveQueue(ctx, args[0], jobclient.SaveQueueOpts{
		Concurrency: 10,
	})
	if err != nil {
		return err
	}
	fmt.Printf("res: %+v\n", q)
	return nil
}
