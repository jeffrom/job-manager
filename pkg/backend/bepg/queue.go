package bepg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/resource"
)

func (pg *Postgres) GetQueue(ctx context.Context, name string) (*resource.Queue, error) {
	c, err := pg.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return getQueueByName(ctx, c, name)
}

func (pg *Postgres) SaveQueue(ctx context.Context, queue *resource.Queue) (*resource.Queue, error) {
	c, err := pg.getConn(ctx)
	if err != nil {
		return nil, err
	}

	prev, err := getQueueByName(ctx, c, queue.Name)
	if err != nil {
		return nil, err
	}

	if prev != nil {
		if !queue.Version.Equals(prev.Version) {
			return nil, backend.NewVersionConflictError(prev.Version, queue.Version, "queue", queue.ID)
		}
		if prev.EqualAttrs(queue) {
			return nil, nil
		}
	}

	queue.Version.Inc()

	q := "INSERT INTO queues (name, v, concurrency, retries, duration, checkin_duration, claim_duration, unique_args, job_schema, created_at) VALUES (:name, :v, :concurrency, :retries, :duration, :checkin_duration, :claim_duration, :unique_args, :job_schema, NOW()) RETURNING *"
	stmt, err := c.PrepareNamedContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	resq := &resource.Queue{}
	if err := stmt.GetContext(ctx, resq, queue); err != nil {
		return nil, err
	}

	return resq, nil
}

func (pg *Postgres) ListQueues(ctx context.Context, opts *resource.QueueListParams) (*resource.Queues, error) {
	return nil, nil
}

func getQueueByName(ctx context.Context, c sqlxer, name string) (*resource.Queue, error) {
	q := "SELECT * FROM queues WHERE name = $1 ORDER BY v DESC LIMIT 1"
	queue := &resource.Queue{}
	if err := sqlx.GetContext(ctx, c, queue, q, name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// TODO populate rest of fields from db (labels etc)
	return queue, nil
}
