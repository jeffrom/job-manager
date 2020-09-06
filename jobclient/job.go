package jobclient

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
)

func (c *Client) EnqueueJob(ctx context.Context, name string, args ...interface{}) (string, error) {
	argList, err := structpb.NewList(args)
	if err != nil {
		return "", err
	}
	params := &apiv1.EnqueueParams{
		Jobs: []*apiv1.EnqueueParamArgs{{Job: name, Args: argList.Values}},
	}

	uri := fmt.Sprintf("/api/v1/queues/%s/enqueue", name)
	req, err := c.newRequestProto("POST", uri, params)
	if err != nil {
		return "", err
	}

	resp := &apiv1.EnqueueResponse{}
	err = c.doRequest(ctx, req, resp)
	if err != nil {
		return "", err
	}

	if len(resp.Jobs) == 0 {
		return "", errors.New("jobclient: unexpectedly received no enqueued job data")
	}

	return resp.Jobs[0], nil
}

func (c *Client) DequeueJobs(ctx context.Context, num int, queueID string, selectors ...string) (*jobv1.Jobs, error) {
	params := &apiv1.DequeueParams{
		Selectors: selectors,
	}
	if num > 0 {
		params.Num = int32(num)
	}
	if queueID != "" {
		params.Job = queueID
	}

	uri := "/api/v1/jobs/dequeue"
	if queueID != "" {
		uri = fmt.Sprintf("/api/v1/queues/%s/dequeue", queueID)
	}
	req, err := c.newRequestProto("POST", uri, params)
	if err != nil {
		return nil, err
	}

	resp := &apiv1.DequeueResponse{}
	err = c.doRequest(ctx, req, resp)
	if err != nil {
		return nil, err
	}
	return resp.Jobs, nil
}

type AckJobOpts struct {
	Data interface{}
}

func (c *Client) AckJob(ctx context.Context, id string, status jobv1.Status) error {
	return c.AckJobOpts(ctx, id, status, AckJobOpts{})
}

func (c *Client) AckJobOpts(ctx context.Context, id string, status jobv1.Status, opts AckJobOpts) error {
	args := &apiv1.AckParamArgs{
		Id:     id,
		Status: status,
	}
	if opts.Data != nil {
		val, err := structpb.NewValue(opts.Data)
		if err != nil {
			return err
		}
		args.Data = val
	}

	uri := "/api/v1/jobs/ack"
	params := &apiv1.AckParams{Acks: []*apiv1.AckParamArgs{args}}
	req, err := c.newRequestProto("POST", uri, params)
	if err != nil {
		return err
	}
	return c.doRequest(ctx, req, nil)
}

func (c *Client) GetJob(ctx context.Context, id string) (*jobv1.Job, error) {
	uri := fmt.Sprintf("/api/v1/jobs/%s", id)
	req, err := c.newRequestProto("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	jobData := &jobv1.Job{}
	if err := c.doRequest(ctx, req, jobData); err != nil {
		return nil, err
	}
	return jobData, nil
}
