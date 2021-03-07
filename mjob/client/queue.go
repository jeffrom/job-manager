package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
	"github.com/jeffrom/job-manager/mjob/resource"
	jobv1 "github.com/jeffrom/job-manager/mjob/resource/job/v1"
	"github.com/jeffrom/job-manager/mjob/schema"
)

type Queue struct {
	*jobv1.Queue
}

type SaveQueueOpts struct {
	MaxRetries      int
	JobDuration     time.Duration
	CheckinDuration time.Duration
	ClaimDuration   time.Duration
	Labels          map[string]string
	Schema          []byte
	Unique          bool
	Version         string
	BackoffInitial  time.Duration
	BackoffMax      time.Duration
	BackoffFactor   float32
}

func (c *Client) SaveQueue(ctx context.Context, name string, opts SaveQueueOpts) (*resource.Queue, error) {
	v, err := resource.NewVersionFromString(opts.Version)
	if err != nil {
		return nil, err
	}
	args := &apiv1.SaveQueueParamArgs{
		Name:   name,
		Labels: opts.Labels,
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
	args.V = v.Raw()
	args.BackoffFactor = opts.BackoffFactor
	if opts.BackoffInitial > 0 {
		args.BackoffInitialDuration = durationpb.New(opts.BackoffInitial)
	}
	if opts.BackoffMax > 0 {
		args.BackoffMaxDuration = durationpb.New(opts.BackoffMax)
	}

	uri := fmt.Sprintf("/api/v1/queues/%s", name)
	req, err := c.newRequestProto(ctx, "PUT", uri, args)
	if err != nil {
		return nil, err
	}

	resp := &apiv1.SaveQueueResponse{}
	if err := c.doRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return jobv1.NewQueueFromProto(resp.Queue), nil
}

type ListQueuesOpts struct {
	Names     []string
	Selectors []string
	Page      *resource.Pagination
}

func (c *Client) ListQueues(ctx context.Context, opts ListQueuesOpts) (*resource.Queues, error) {
	page := apiv1.PaginationToProto(opts.Page)
	params := &apiv1.ListQueuesRequest{
		Names:     opts.Names,
		Selectors: opts.Selectors,
		Page:      page,
	}
	uri := "/api/v1/queues"
	req, err := c.newRequestProto(ctx, "GET", uri, params)
	if err != nil {
		return nil, err
	}

	resp := &apiv1.ListQueuesResponse{}
	if err := c.doRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return &resource.Queues{Queues: jobv1.NewQueuesFromProto(resp.Items)}, nil
}

func (c *Client) GetQueue(ctx context.Context, id string) (*resource.Queue, error) {
	uri := fmt.Sprintf("/api/v1/queues/%s", id)
	req, err := c.newRequestProto(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}

	resp := &apiv1.GetQueueResponse{}
	if err := c.doRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return jobv1.NewQueueFromProto(resp.Data), nil
}
