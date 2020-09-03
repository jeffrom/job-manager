package jobclient

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/job"
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

	if len(resp.Jobs.Jobs) == 0 {
		return "", errors.New("jobclient: unexpectedly received no enqueued job data")
	}

	return resp.Jobs.Jobs[0].Id, nil
}

func (c *Client) DequeueJobs(ctx context.Context, num int, queueID string, selectors ...string) (*job.Jobs, error) {
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
	Data map[string]interface{}
}

func (c *Client) AckJob(ctx context.Context, id string, status job.Status) error {
	return c.AckJobOpts(ctx, id, status, AckJobOpts{})
}

func (c *Client) AckJobOpts(ctx context.Context, id string, status job.Status, opts AckJobOpts) error {
	args := &apiv1.AckParamArgs{
		Id:     id,
		Status: status,
	}
	if len(opts.Data) > 0 {
		data, err := structpb.NewStruct(opts.Data)
		if err != nil {
			return err
		}
		args.Data = data.Fields
	}

	uri := "/api/v1/jobs/ack"
	params := &apiv1.AckParams{Acks: []*apiv1.AckParamArgs{args}}
	req, err := c.newRequestProto("POST", uri, params)
	if err != nil {
		return err
	}
	return c.doRequest(ctx, req, nil)
}

func (c *Client) GetJob(ctx context.Context, id string) (*job.Job, error) {
	uri := fmt.Sprintf("/api/v1/jobs/%s", id)
	req, err := c.newRequestProto("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	jobData := &job.Job{}
	if err := c.doRequest(ctx, req, jobData); err != nil {
		return nil, err
	}
	return jobData, nil
}
