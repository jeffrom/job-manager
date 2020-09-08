package jobclient

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/label"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
)

type EnqueueOpts struct {
	Data   interface{}
	Claims label.Claims
}

func (c *Client) EnqueueJobOpts(ctx context.Context, name string, opts EnqueueOpts, args ...interface{}) (string, error) {
	argList, err := structpb.NewList(args)
	if err != nil {
		return "", err
	}

	var jobData *jobv1.Data
	if opts.Data != nil || len(opts.Claims) > 0 {
		jobData = &jobv1.Data{}
		if opts.Data != nil {
			data, err := structpb.NewValue(opts.Data)
			if err != nil {
				return "", err
			}
			jobData.Data = data
		}
		if len(opts.Claims) > 0 {
			jobData.Claims = opts.Claims.Format()
		}
	}
	params := &apiv1.EnqueueRequest{
		Jobs: []*apiv1.EnqueueRequestArgs{
			{
				Job:  name,
				Args: argList.Values,
				Data: jobData,
			},
		},
	}

	uri := fmt.Sprintf("/api/v1/queues/%s/enqueue", name)
	req, err := c.newRequestProto(ctx, "POST", uri, params)
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

func (c *Client) EnqueueJob(ctx context.Context, name string, args ...interface{}) (string, error) {
	return c.EnqueueJobOpts(ctx, name, EnqueueOpts{}, args...)
}

type DequeueOpts struct {
	Queues    []string
	Selectors []string
	Claims    label.Claims
}

func (c *Client) DequeueJobsOpts(ctx context.Context, num int, opts DequeueOpts) (*jobv1.Jobs, error) {
	queueID := ""
	if len(opts.Queues) == 1 {
		queueID = opts.Queues[0]
	}
	params := &apiv1.DequeueRequest{
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

	resp := &apiv1.DequeueResponse{}
	err = c.doRequest(ctx, req, resp)
	if err != nil {
		return nil, err
	}
	return resp.Jobs, nil
}

func (c *Client) DequeueJobs(ctx context.Context, num int, queueID string) (*jobv1.Jobs, error) {
	return c.DequeueJobsOpts(ctx, num, DequeueOpts{Queues: []string{queueID}})
}

type AckJobOpts struct {
	Data interface{}
}

func (c *Client) AckJob(ctx context.Context, id string, status jobv1.Status) error {
	return c.AckJobOpts(ctx, id, status, AckJobOpts{})
}

func (c *Client) AckJobOpts(ctx context.Context, id string, status jobv1.Status, opts AckJobOpts) error {
	args := &apiv1.AckRequestArgs{
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
	params := &apiv1.AckRequest{Acks: []*apiv1.AckRequestArgs{args}}
	req, err := c.newRequestProto(ctx, "POST", uri, params)
	if err != nil {
		return err
	}
	return c.doRequest(ctx, req, nil)
}

func (c *Client) GetJob(ctx context.Context, id string) (*jobv1.Job, error) {
	uri := fmt.Sprintf("/api/v1/jobs/%s", id)
	req, err := c.newRequestProto(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}

	jobData := &jobv1.Job{}
	if err := c.doRequest(ctx, req, jobData); err != nil {
		return nil, err
	}
	return jobData, nil
}
