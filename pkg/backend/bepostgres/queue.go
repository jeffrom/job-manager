package bepostgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/jeffrom/job-manager/mjob/label"
	"github.com/jeffrom/job-manager/mjob/resource"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/internal"
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
	if queue.Version == nil {
		queue.Version = resource.NewVersion(0)
	}

	prev, err := getQueueByName(ctx, c, queue.Name)
	if err != nil && err != backend.ErrNotFound {
		return nil, err
	}

	// fmt.Printf("prev: %+v\n", prev)
	// fmt.Printf("curr: %+v\n", queue)
	if prev != nil {
		if !queue.Version.Equals(prev.Version) {
			return nil, backend.NewVersionConflictError(prev.Version, queue.Version, "queue", queue.Name)
		}
		// fmt.Println("equal:", prev.EqualAttrs(queue))
		if prev.EqualAttrs(queue) {
			return prev, nil
		}
	}

	now := internal.GetTimeProvider(ctx).Now().UTC()
	if queue.CreatedAt.IsZero() {
		queue.CreatedAt = now
	}
	queue.UpdatedAt = now
	queue.Version.Inc()

	results, err := insertQueues(ctx, c, []*resource.Queue{queue})
	if err != nil {
		return nil, err
	}
	return results[0], nil
}

func (pg *Postgres) ListQueues(ctx context.Context, opts *resource.QueueListParams) (*resource.Queues, error) {
	c, err := pg.getConn(ctx)
	if err != nil {
		return nil, err
	}

	var wheres []string
	var joins []string
	var args []interface{}
	q := "SELECT DISTINCT ON (name) id, queues.name, v, concurrency, retries, unique_args, duration, checkin_duration, claim_duration, job_schema, created_at FROM queues"
	if opts != nil {
		if len(opts.Names) > 0 {
			wheres = append(wheres, "name IN (?)")
			args = append(args, opts.Names)
		}

		joins, wheres, args = sqlSelectors(opts.Selectors, joins, wheres, args)
	}
	if len(joins) > 0 {
		q += " " + strings.Join(joins, " ")
	}
	if len(wheres) > 0 {
		q += " WHERE (" + strings.Join(wheres, " AND ") + ")"
	}
	q += " ORDER BY name, v DESC"
	q, args, err = sqlx.In(q, args...)
	if err != nil {
		return nil, err
	}

	rows := []*resource.Queue{}
	if err := sqlx.SelectContext(ctx, c, &rows, c.Rebind(q), args...); err != nil {
		return nil, err
	}

	if err := annotateQueues(ctx, c, rows); err != nil {
		return nil, err
	}

	if opts != nil && opts.Selectors.Len() > 0 {
		var frows []*resource.Queue
		for _, row := range rows {
			if !opts.Selectors.Match(row.Labels) {
				continue
			}
			frows = append(frows, row)
		}
		rows = frows
	}

	return &resource.Queues{Queues: rows}, nil
}

func getQueueByName(ctx context.Context, c sqlxer, name string) (*resource.Queue, error) {
	q := "SELECT * FROM queues WHERE name = $1 ORDER BY v DESC LIMIT 1"
	queue := &resource.Queue{}
	if err := sqlx.GetContext(ctx, c, queue, q, name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, backend.ErrNotFound
		}
		return nil, err
	}

	if err := annotateQueues(ctx, c, []*resource.Queue{queue}); err != nil {
		return nil, err
	}
	return queue, nil
}

type queueLabelRow struct {
	Queue string
	Name  string
	Value string
}

func annotateQueues(ctx context.Context, c sqlxer, queues []*resource.Queue) error {
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

	q := "INSERT INTO queues (name, v, concurrency, retries, duration, checkin_duration, claim_duration, unique_args, job_schema, created_at, updated_at) VALUES (:name, :v, :concurrency, :retries, :duration, :checkin_duration, :claim_duration, :unique_args, :job_schema, :created_at, :updated_at) RETURNING *"
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

	labelDelQ := "DELETE FROM queue_labels WHERE queue = ? AND name NOT IN (?)"
	// TODO labelDelAllQ for when len(queue.Labels) == 0 && queue.Labels != nil?

	results := make([]*resource.Queue, len(queues))
	for i, queue := range queues {
		resq := &resource.Queue{}
		if err := stmt.GetContext(ctx, resq, queue); err != nil {
			return nil, err
		}

		names := make([]string, len(queue.Labels))
		ii := 0
		for name, val := range queue.Labels {
			if _, err := labelstmt.ExecContext(ctx, queue.Name, name, val); err != nil {
				return nil, err
			}
			names[ii] = name
			ii++
		}

		// fmt.Println(names)
		// delete labels if they are modified
		if len(names) > 0 {
			lq, labelDelArgs, err := sqlx.In(labelDelQ, queue.Name, names)
			if err != nil {
				return nil, err
			}
			if _, err := c.ExecContext(ctx, c.Rebind(lq), labelDelArgs...); err != nil {
				return nil, err
			}
		} else if names != nil {
			// TODO del all labels
		}

		resq.Labels = queue.Labels

		results[i] = resq
	}
	return results, nil
}
