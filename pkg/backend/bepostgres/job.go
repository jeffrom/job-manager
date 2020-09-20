package bepostgres

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"strconv"

	"github.com/jeffrom/job-manager/pkg/internal"
	"github.com/jeffrom/job-manager/pkg/resource"
	"github.com/jmoiron/sqlx"
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

	// var uniquenessKeys []string
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
	q := "INSERT INTO jobs (" + fields + ", root_id) VALUES (" + vals + ", 0) RETURNING id"
	stmt, err := c.PrepareNamedContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for _, jb := range jobs.Jobs {
		row := stmt.QueryRowContext(ctx, jb)
		var id int64
		if err := row.Scan(&id); err != nil {
			return nil, err
		}
		jb.ID = strconv.FormatInt(id, 10)
	}

	return jobs, nil
}

func (pg *Postgres) DequeueJobs(ctx context.Context, limit int, opts *resource.JobListParams) (*resource.Jobs, error) {
	return nil, nil
}

func (pg *Postgres) AckJobs(ctx context.Context, results *resource.Acks) error {
	return nil
}

func (pg *Postgres) GetJobByID(ctx context.Context, id string) (*resource.Job, error) {
	// id, err := strconv.ParseInt(idstr, 10, 64)
	// if err != nil {
	// 	return nil, err
	// }
	c, err := pg.getConn(ctx)
	if err != nil {
		return nil, err
	}

	jb := &resource.Job{}
	if err := sqlx.GetContext(ctx, c, jb, "SELECT * FROM jobs WHERE id = $1", id); err != nil {
		return nil, err
	}
	return jb, nil
}

func (pg *Postgres) ListJobs(ctx context.Context, limit int, opts *resource.JobListParams) (*resource.Jobs, error) {
	return nil, nil
}

func uniquenessKeyFromArgs(args []interface{}) (string, error) {
	b, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return string(sum[:]), nil
}
