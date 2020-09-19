package bepg

import (
	"context"

	"github.com/jeffrom/job-manager/pkg/internal"
	"github.com/jeffrom/job-manager/pkg/resource"
)

func (pg *Postgres) EnqueueJobs(ctx context.Context, jobs *resource.Jobs) (*resource.Jobs, error) {
	c, err := pg.getConn(ctx)
	if err != nil {
		return nil, err
	}
	res, err := pg.ListQueues(ctx, &resource.QueueListParams{Names: jobs.Queues()})
	if err != nil {
		return nil, err
	}
	queues := res.ToMap()

	now := internal.GetTimeProvider(ctx).Now().UTC()
	for _, jb := range jobs.Jobs {
		q := queues[jb.Name]
		jb.EnqueuedAt = now
		jb.QueueID = q.ID
		jb.Version = resource.NewVersion(1)
		jb.QueueVersion = q.Version
		jb.Status = resource.NewStatus(resource.StatusQueued)
	}

	fields, vals := sqlFields(
		"v", "queue", "queue_v",
		"attempt", "status",
		"args", "data", "enqueued_at",
	)
	q := "INSERT INTO jobs " + fields + " VALUES " + vals
	stmt, err := c.PrepareNamedContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for _, jb := range jobs.Jobs {
		if _, err := stmt.ExecContext(ctx, jb); err != nil {
			return nil, err
		}
	}

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
