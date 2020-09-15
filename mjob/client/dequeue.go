package client

import (
	"context"
	"fmt"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/label"
	"github.com/jeffrom/job-manager/pkg/resource"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
)

type DequeueOpts struct {
	Queues    []string
	Selectors []string
	Claims    label.Claims
}

func (c *Client) DequeueJobsOpts(ctx context.Context, num int, opts DequeueOpts) (*resource.Jobs, error) {
	queueID := ""
	if len(opts.Queues) == 1 {
		queueID = opts.Queues[0]
	}
	params := &apiv1.DequeueJobsRequest{
		Queues:    opts.Queues,
		Selectors: opts.Selectors,
		Claims:    opts.Claims.Format(),
	}
	if num > 0 {
		params.Num = int32(num)
	}

	uri := "/api/v1/jobs/dequeue"
	if queueID != "" {
		uri = fmt.Sprintf("/api/v1/queues/%s/dequeue", queueID)
	}
	req, err := c.newRequestProto(ctx, "POST", uri, params)
	if err != nil {
		return nil, err
	}

	resp := &apiv1.DequeueJobsResponse{}
	err = c.doRequest(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	resJobs, err := jobv1.NewJobsFromProto(resp.Jobs.Jobs)
	if err != nil {
		return nil, err
	}
	return &resource.Jobs{Jobs: resJobs}, nil
}

func (c *Client) DequeueJobs(ctx context.Context, num int, queueID string) (*resource.Jobs, error) {
	return c.DequeueJobsOpts(ctx, num, DequeueOpts{Queues: []string{queueID}})
}
