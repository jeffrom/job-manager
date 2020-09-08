package client

import (
	"context"

	"google.golang.org/protobuf/types/known/structpb"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/resource"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
)

type AckJobOpts struct {
	Data interface{}
}

func (c *Client) AckJob(ctx context.Context, id string, status resource.Status) error {
	return c.AckJobOpts(ctx, id, status, AckJobOpts{})
}

func (c *Client) AckJobOpts(ctx context.Context, id string, status resource.Status, opts AckJobOpts) error {
	args := &apiv1.AckJobsRequestArgs{
		Id:     id,
		Status: jobv1.Status(status),
	}
	if opts.Data != nil {
		val, err := structpb.NewValue(opts.Data)
		if err != nil {
			return err
		}
		args.Data = val
	}

	uri := "/api/v1/jobs/ack"
	params := &apiv1.AckJobsRequest{Acks: []*apiv1.AckJobsRequestArgs{args}}
	req, err := c.newRequestProto(ctx, "POST", uri, params)
	if err != nil {
		return err
	}
	return c.doRequest(ctx, req, nil)
}
