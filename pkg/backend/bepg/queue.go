package bepg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/internal"
	"github.com/jeffrom/job-manager/pkg/label"
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

	fmt.Printf("prev: %+v\n", prev)
	fmt.Printf("curr: %+v\n", queue)
	if prev != nil {
		if !queue.Version.Equals(prev.Version) {
			return nil, backend.NewVersionConflictError(prev.Version, queue.Version, "queue", queue.ID)
		}
		fmt.Println("equal:", prev.EqualAttrs(queue))
		if prev.EqualAttrs(queue) {
			return nil, nil
		}
	}

	now := internal.GetTimeProvider(ctx).Now()
	queue.CreatedAt = now
	queue.Version.Inc()

	results, err := insertQueues(ctx, c, []*resource.Queue{queue})
	if err != nil {
		return nil, err
	}
	return results[0], nil
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

	if err := annotateQueuesWithLabels(ctx, c, []*resource.Queue{queue}); err != nil {
		return nil, err
	}
	return queue, nil
}

type queueLabelRow struct {
	Queue string
	Name  string
	Value string
}

func annotateQueuesWithLabels(ctx context.Context, c sqlxer, queues []*resource.Queue) error {
	if len(queues) == 0 {
		return nil
	}

	names := make([]string, len(queues))
	for i, queue := range queues {
		names[i] = queue.Name
	}

	q, args, err := sqlx.In("SELECT * FROM queue_labels WHERE queue IN (?)", names)
	if err != nil {
		return err
	}

	rows := []*queueLabelRow{}
	if err := sqlx.SelectContext(ctx, c, &rows, c.Rebind(q), args...); err != nil {
		return err
	}

	// build labels
	labelmap := make(map[string]label.Labels)
	for _, row := range rows {
		// fmt.Printf("queue_labels row: %q %q %q\n", row.Queue, row.Name, row.Value)
		labels, ok := labelmap[row.Queue]
		if !ok {
			labels = make(label.Labels)
		}

		labels[row.Name] = row.Value

		labelmap[row.Queue] = labels
	}
	for _, queue := range queues {
		queue.Labels = labelmap[queue.Name]
	}
	return nil
}

func insertQueues(ctx context.Context, c sqlxer, queues []*resource.Queue) ([]*resource.Queue, error) {
	if len(queues) == 0 {
		return nil, nil
	}

	q := "INSERT INTO queues (name, v, concurrency, retries, duration, checkin_duration, claim_duration, unique_args, job_schema, created_at) VALUES (:name, :v, :concurrency, :retries, :duration, :checkin_duration, :claim_duration, :unique_args, :job_schema, :created_at) RETURNING *"
	stmt, err := c.PrepareNamedContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	labelq := "INSERT INTO queue_labels (queue, name, value) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING"
	labelstmt, err := c.PrepareContext(ctx, labelq)
	if err != nil {
		return nil, err
	}
	defer labelstmt.Close()

	results := make([]*resource.Queue, len(queues))
	for i, queue := range queues {
		resq := &resource.Queue{}
		if err := stmt.GetContext(ctx, resq, queue); err != nil {
			return nil, err
		}

		for name, val := range queue.Labels {
			if _, err := labelstmt.ExecContext(ctx, queue.Name, name, val); err != nil {
				return nil, err
			}
		}
		resq.Labels = queue.Labels

		results[i] = resq
	}
	return results, nil
}
