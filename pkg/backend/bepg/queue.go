package bepg

import (
	"context"

	"github.com/jmoiron/sqlx"

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
	if prev.EqualAttrs(queue) {
		return nil, nil
	}

	// q := "INSERT INTO queues"
	return nil, nil
}

func (pg *Postgres) ListQueues(ctx context.Context, opts *resource.QueueListParams) (*resource.Queues, error) {
	return nil, nil
}

func getQueueByID(ctx context.Context, c sqlx.ExtContext, id string) (*resource.Queue, error) {
	q := "SELECT * FROM queues WHERE id = ?"
	queue := &resource.Queue{}
	if err := sqlx.GetContext(ctx, c, q, id); err != nil {
		return nil, err
	}

	// TODO populate rest of fields from db (labels etc)
	return queue, nil
}
