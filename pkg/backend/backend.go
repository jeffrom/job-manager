// Package backend manages storing data.
package backend

import (
	"context"

	"github.com/jeffrom/job-manager/pkg/job"
)

type Interface interface {
	GetQueue(ctx context.Context, job string) (*job.Queue, error)
	SaveQueue(ctx context.Context, queue *job.Queue) error
	ListQueues(ctx context.Context, opts *job.QueueListParams) (*job.Queues, error)

	EnqueueJobs(ctx context.Context, jobs *job.Jobs) error
	DequeueJobs(ctx context.Context, num int, opts *job.JobListParams) (*job.Jobs, error)
	AckJobs(ctx context.Context, results *job.Acks) error
	GetSetJobKeys(ctx context.Context, keys []string) (bool, error)
	DeleteJobKeys(ctx context.Context, keys []string) error

	GetJobByID(ctx context.Context, id string) (*job.Job, error)
	ListJobs(ctx context.Context, opts *job.JobListParams) (*job.Jobs, error)
}
