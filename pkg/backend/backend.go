// Package backend manages storing data.
package backend

import (
	"context"

	"github.com/jeffrom/job-manager/pkg/resource"
)

// Interface defines the backend interface. Some required properties:
//
// - UTC
//
// - handle resource versions, conflicts
//
// - safe for concurrent operations
type Interface interface {
	Ping(ctx context.Context) error
	Reset(ctx context.Context) error

	GetQueue(ctx context.Context, job string) (*resource.Queue, error)
	SaveQueue(ctx context.Context, queue *resource.Queue) (*resource.Queue, error)
	ListQueues(ctx context.Context, opts *resource.QueueListParams) (*resource.Queues, error)

	EnqueueJobs(ctx context.Context, jobs *resource.Jobs) (*resource.Jobs, error)
	DequeueJobs(ctx context.Context, num int, opts *resource.JobListParams) (*resource.Jobs, error)
	AckJobs(ctx context.Context, results *resource.Acks) error
	GetSetJobKeys(ctx context.Context, keys []string) (bool, error)
	DeleteJobKeys(ctx context.Context, keys []string) error

	GetJobByID(ctx context.Context, id string) (*resource.Job, error)
	ListJobs(ctx context.Context, opts *resource.JobListParams) (*resource.Jobs, error)
}
