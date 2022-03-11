package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/label"
	"github.com/jeffrom/job-manager/mjob/resource"
)

type enqueueOpts struct {
	// client.EnqueueOpts
	data   string
	claims []string
}

type enqueueCmd struct {
	*cobra.Command
	opts *enqueueOpts
}

func newEnqueueCmd(cfg *client.Config) *enqueueCmd {
	opts := &enqueueOpts{}
	c := &enqueueCmd{
		Command: &cobra.Command{
			Use:               "enqueue",
			Short:             "Enqueue a job",
			Aliases:           []string{"enq"},
			Args:              cobra.MinimumNArgs(1),
			ValidArgsFunction: validQueueList(1),
			Long:              `Enqueue a job.`,
			Example: `# enqueue a job with args (in json format) [1, 2, 3]
$ jobctl enqueue myq 1 2 3

# enqueue a job with args [{"hi": "hello"}]
$ jobctl enqueue myq '{"hi": "hello"}'

# enqueue a job with a claim
$ jobctl enqueue -c myclaim=myval myq myarg

# enqueue with some data
$ jobctl enqueue myq myarg --data-raw '{"hi": "hello"}'
`,
		},
		opts: opts,
	}

	flags := c.Command.Flags()
	flags.StringArrayVarP(&opts.claims, "claim", "c", nil, "enqueue with claims in `key=val` format. Multiple values can be provided")
	flags.StringVar(&opts.data, "data-raw", "", "raw `json` data to include with job")

	return c
}

func (c *enqueueCmd) Cmd() *cobra.Command { return c.Command }
func (c *enqueueCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	var claims label.Claims
	if len(c.opts.claims) > 0 {
		var err error
		claims, err = label.ParseClaims(c.opts.claims)
		if err != nil {
			return err
		}
	}

	var data interface{}
	dataRaw := c.opts.data
	if dataRaw != "" {
		if err := json.Unmarshal([]byte(dataRaw), &data); err != nil {
			return err
		}
	}
	iargs, err := resource.ParseCLIArgs(args[1:])
	if err != nil {
		return err
	}
	opts := client.EnqueueOpts{
		Claims: claims,
		Data:   data,
	}

	cl := clientFromContext(ctx)
	id, err := cl.EnqueueJobOpts(ctx, args[0], opts, iargs...)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", id)
	return nil
}
