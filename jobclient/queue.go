package jobclient

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/job"
)

type SaveQueueOptions struct {
	Concurrency int
	MaxRetries  int
	JobDuration time.Duration
	Labels      map[string]string
}

func (c *Client) SaveQueue(ctx context.Context, params string, opts SaveQueueOptions) (*job.Queue, error) {
	args := &apiv1.SaveQueueParamArgs{
		Name:   params,
		Labels: opts.Labels,
	}
	if opts.Concurrency > 0 {
		args.Concurrency = proto.Int32(int32(opts.Concurrency))
	}
	if opts.MaxRetries > 0 {
		args.MaxRetries = proto.Int32(int32(opts.MaxRetries))
	}
	if opts.JobDuration > 0 {
		args.Duration = durationpb.New(opts.JobDuration)
	}

	uri := fmt.Sprintf("/api/v1/jobs/%s", params)
	req, err := c.newRequestProto("PUT", uri, args)
	if err != nil {
		return nil, err
	}

	queue := &job.Queue{}
	if err := c.doRequest(ctx, req, queue); err != nil {
		return nil, err
	}
	return queue, nil
}
