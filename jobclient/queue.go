package jobclient

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
	"github.com/jeffrom/job-manager/pkg/schema"
)

type Queue struct {
	*jobv1.Queue
}

type SaveQueueOpts struct {
	Concurrency     int
	MaxRetries      int
	JobDuration     time.Duration
	CheckinDuration time.Duration
	ClaimDuration   time.Duration
	Labels          map[string]string
	Schema          []byte
	ArgSchema       []byte
	DataSchema      []byte
	ResultSchema    []byte
	Unique          bool
	V               int32
}

func (c *Client) SaveQueue(ctx context.Context, name string, opts SaveQueueOpts) (*jobv1.Queue, error) {
	args := &apiv1.SaveQueueParamArgs{
		Name:   name,
		Labels: opts.Labels,
	}
	if opts.Concurrency > 0 {
		args.Concurrency = int32(opts.Concurrency)
	}
	if opts.MaxRetries > 0 {
		args.MaxRetries = int32(opts.MaxRetries)
	}
	if opts.JobDuration > 0 {
		args.Duration = durationpb.New(opts.JobDuration)
	}
	if opts.ClaimDuration > 0 {
		args.ClaimDuration = durationpb.New(opts.ClaimDuration)
	}
	if opts.CheckinDuration > 0 {
		args.CheckinDuration = durationpb.New(opts.CheckinDuration)
	}
	if len(opts.Schema) > 0 {
		cSchema, err := schema.Canonicalize(opts.Schema)
		if err != nil {
			return nil, err
		}
		args.Schema = cSchema
	}
	args.Unique = opts.Unique
	args.V = opts.V

	uri := fmt.Sprintf("/api/v1/queues/%s", name)
	req, err := c.newRequestProto("PUT", uri, args)
	if err != nil {
		return nil, err
	}

	resp := &apiv1.SaveQueueResponse{}
	if err := c.doRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp.Queue, nil
}

type ListQueuesOpts struct {
	Names     []string
	Selectors []string
}

func (c *Client) ListQueues(ctx context.Context, opts ListQueuesOpts) (*jobv1.Queues, error) {
	params := &apiv1.ListQueuesRequest{
		Names:     opts.Names,
		Selectors: opts.Selectors,
	}
	uri := "/api/v1/queues"
	req, err := c.newRequestProto("GET", uri, params)
	if err != nil {
		return nil, err
	}

	resp := &apiv1.ListQueuesResponse{}
	if err := c.doRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) GetQueue(ctx context.Context, id string) (*jobv1.Queue, error) {
	uri := fmt.Sprintf("/api/v1/queues/%s", id)
	req, err := c.newRequestProto("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	resp := &apiv1.GetQueueResponse{}
	if err := c.doRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
