package bepg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/resource"
)

func (pg *Postgres) GetQueue(ctx context.Context, id string) (*resource.Queue, error) {
	c, err := pg.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return getQueueByID(ctx, c, id)
}

func (pg *Postgres) SaveQueue(ctx context.Context, queue *resource.Queue) (*resource.Queue, error) {
	c, err := pg.getConn(ctx)
	if err != nil {
		return nil, err
	}

	prev, err := getQueueByID(ctx, c, queue.ID)
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

	// q := "INSERT INTO queues"
	return nil, nil
}

func (pg *Postgres) ListQueues(ctx context.Context, opts *resource.QueueListParams) (*resource.Queues, error) {
	return nil, nil
}

func getQueueByID(ctx context.Context, c sqlx.ExtContext, id string) (*resource.Queue, error) {
	q := "SELECT * FROM queues WHERE id = $1"
	queue := &resource.Queue{}
	if err := sqlx.GetContext(ctx, c, queue, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// TODO populate rest of fields from db (labels etc)
	return queue, nil
}
