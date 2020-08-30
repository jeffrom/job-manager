package jobclient

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"
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

	uri := fmt.Sprintf("/api/v1/jobs/%s/enqueue", name)
	req, err := c.newRequestProto("POST", uri, params)
	if err != nil {
		return "", err
	}

	jobs := &job.Jobs{}
	err = c.doRequest(ctx, req, jobs)
	if err != nil {
		return "", err
	}

	return jobs.Jobs[0].Id, nil
}

func (c *Client) DequeueJobs(ctx context.Context, num int, jobName string, selectors ...string) (*job.Jobs, error) {
	params := &apiv1.DequeueParams{
		Selectors: selectors,
	}
	if num > 0 {
		params.Num = proto.Int32(int32(num))
	}
	if jobName != "" {
		params.Job = proto.String(jobName)
	}

	uri := "/api/v1/jobs/dequeue"
	if jobName != "" {
		uri = fmt.Sprintf("/api/v1/jobs/%s/dequeue", jobName)
	}
	req, err := c.newRequestProto("POST", uri, params)
	if err != nil {
		return nil, err
	}

	jobs := &job.Jobs{}
	err = c.doRequest(ctx, req, jobs)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (c *Client) AckJob(ctx context.Context, status job.Status) error {
	return nil
}
