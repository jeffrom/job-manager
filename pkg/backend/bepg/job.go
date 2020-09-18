package bepg

import (
	"context"

	"github.com/jeffrom/job-manager/pkg/resource"
)

func (pg *Postgres) EnqueueJobs(ctx context.Context, jobs *resource.Jobs) (*resource.Jobs, error) {
	return nil, nil
}

func (pg *Postgres) DequeueJobs(ctx context.Context, limit int, opts *resource.JobListParams) (*resource.Jobs, error) {
	return nil, nil
}

func (pg *Postgres) AckJobs(ctx context.Context, results *resource.Acks) error {
	return nil
}

func (pg *Postgres) GetJobByID(ctx context.Context, id string) (*resource.Job, error) {
	return nil, nil
}

func (pg *Postgres) ListJobs(ctx context.Context, limit int, opts *resource.JobListParams) (*resource.Jobs, error) {
	return nil, nil
}
