package client

import (
	"context"
	"encoding/json"

	apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
	"github.com/jeffrom/job-manager/mjob/resource"
	jobv1 "github.com/jeffrom/job-manager/mjob/resource/job/v1"
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
		b, err := json.Marshal(opts.Data)
		if err != nil {
			return err
		}
		args.Data = b
	}

	uri := "/api/v1/jobs/ack"
	params := &apiv1.AckJobsRequest{Acks: []*apiv1.AckJobsRequestArgs{args}}
	req, err := c.newRequestProto(ctx, "POST", uri, params)
	if err != nil {
		return err
	}
	return c.doRequest(ctx, req, nil)
}
