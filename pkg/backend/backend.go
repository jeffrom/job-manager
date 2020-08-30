// Package backend manages storing data.
package backend

import (
	"context"
	"errors"

	"github.com/jeffrom/job-manager/pkg/job"
)

type Interface interface {
	GetQueue(ctx context.Context, job string) (*job.Queue, error)
	SaveQueue(ctx context.Context, queue *job.Queue) error
	ListQueues(ctx context.Context, opts *job.ListOpts) (*job.Queues, error)

	EnqueueJobs(ctx context.Context, jobs *job.Jobs) error
	DequeueJobs(ctx context.Context, num int, opts *job.ListOpts) (*job.Jobs, error)
	AckJobs(ctx context.Context, results *job.Results) error

	GetJobByID(ctx context.Context, id string) (*job.Job, error)
	ListJobs(ctx context.Context, opts *job.ListOpts) (*job.Jobs, error)
}

var ErrNotFound = errors.New("backend: not found")
