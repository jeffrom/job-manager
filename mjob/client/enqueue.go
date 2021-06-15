package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
	"github.com/jeffrom/job-manager/mjob/label"
	jobv1 "github.com/jeffrom/job-manager/mjob/resource/job/v1"
)

type EnqueueOpts struct {
	Data   interface{}
	Claims label.Claims
}

func (c *Client) EnqueueJobOpts(ctx context.Context, name string, opts EnqueueOpts, args ...interface{}) (string, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return "", err
	}

	var jobData *jobv1.Data
	if opts.Data != nil || len(opts.Claims) > 0 {
		jobData = &jobv1.Data{}
		if opts.Data != nil {
			b, err := json.Marshal(opts.Data)
			if err != nil {
				return "", err
			}
			jobData.Data = b
		}
		if len(opts.Claims) > 0 {
			jobData.Claims = opts.Claims.Format()
		}
	}
	params := &apiv1.EnqueueJobsRequest{
		Jobs: []*apiv1.EnqueueJobsRequestArgs{
			{
				Job:  name,
				Args: argsJSON,
				Data: jobData,
			},
		},
	}

	uri := fmt.Sprintf("/api/v1/queues/%s/enqueue", name)
	req, err := c.newRequestProto(ctx, "POST", uri, params)
	if err != nil {
		return "", err
	}

	resp := &apiv1.EnqueueJobsResponse{}
	err = c.doRequest(ctx, req, resp)
	if err != nil {
		return "", err
	}

	if len(resp.Jobs) == 0 {
		return "", errors.New("client: unexpectedly received no enqueued job data")
	}

	return resp.Jobs[0], nil
}

func (c *Client) EnqueueJob(ctx context.Context, name string, args ...interface{}) (string, error) {
	return c.EnqueueJobOpts(ctx, name, EnqueueOpts{}, args...)
}
