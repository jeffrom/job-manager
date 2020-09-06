// Package backend manages storing data.
package backend

import (
	"context"

	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
)

type Interface interface {
	GetQueue(ctx context.Context, job string) (*jobv1.Queue, error)
	SaveQueue(ctx context.Context, queue *jobv1.Queue) error
	ListQueues(ctx context.Context, opts *jobv1.QueueListParams) (*jobv1.Queues, error)

	EnqueueJobs(ctx context.Context, jobs *jobv1.Jobs) error
	DequeueJobs(ctx context.Context, num int, opts *jobv1.JobListParams) (*jobv1.Jobs, error)
	AckJobs(ctx context.Context, results *jobv1.Acks) error
	GetSetJobKeys(ctx context.Context, keys []string) (bool, error)
	DeleteJobKeys(ctx context.Context, keys []string) error

	GetJobByID(ctx context.Context, id string) (*jobv1.Job, error)
	ListJobs(ctx context.Context, opts *jobv1.JobListParams) (*jobv1.Jobs, error)
}
