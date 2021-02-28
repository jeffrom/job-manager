// Package backend manages storing data.
package backend

import (
	"context"
	"net/http"

	"github.com/jeffrom/job-manager/mjob/resource"
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
	// Reset resets the backend. For testing.
	Reset(ctx context.Context) error

	GetQueue(ctx context.Context, job string) (*resource.Queue, error)
	SaveQueue(ctx context.Context, queue *resource.Queue) (*resource.Queue, error)
	ListQueues(ctx context.Context, opts *resource.QueueListParams) (*resource.Queues, error)

	EnqueueJobs(ctx context.Context, jobs *resource.Jobs) (*resource.Jobs, error)
	DequeueJobs(ctx context.Context, limit int, opts *resource.JobListParams) (*resource.Jobs, error)
	AckJobs(ctx context.Context, results *resource.Acks) error

	GetJobUniqueArgs(ctx context.Context, keys []string) ([]string, bool, error)
	SetJobUniqueArgs(ctx context.Context, ids, keys []string) error
	DeleteJobUniqueArgs(ctx context.Context, ids, keys []string) error
	// GetSetJobKeys(ctx context.Context, ids, keys []string) (string, bool, error)
	// DeleteJobKeys(ctx context.Context, keys []string) error

	GetJobByID(ctx context.Context, id string) (*resource.Job, error)
	ListJobs(ctx context.Context, limit int, opts *resource.JobListParams) (*resource.Jobs, error)

	Stats(ctx context.Context, queue string) (*resource.Stats, error)
}

type MiddlewareProvider interface {
	Middleware() func(next http.Handler) http.Handler
}

type HandlerProvider interface {
	Handler() http.Handler
}
