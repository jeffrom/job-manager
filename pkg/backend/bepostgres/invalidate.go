package bepostgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/jeffrom/job-manager/pkg/internal"
)

func (pg *Postgres) InvalidateJobs(ctx context.Context) error {
	if err := pg.ensureConn(ctx); err != nil {
		return err
	}

	c := pg.db
	sql := "UPDATE jobs SET status = 'failed', v=v+1 FROM (SELECT jobs.id FROM jobs LEFT JOIN queues ON jobs.queue_id = queues.id WHERE status = 'running' AND jobs.started_at + (queues.duration / 1000) * interval '1 microsecond' < ? ORDER BY jobs.id ASC LIMIT 50 FOR UPDATE OF jobs SKIP LOCKED) AS to_update WHERE jobs.id = to_update.id"
	q, err := sqlx.PreparexContext(ctx, c, c.Rebind(sql))
	if err != nil {
		return err
	}
	defer q.Close()

	now := internal.GetTimeProvider(ctx).Now().UTC()
	for {
		res, err := q.ExecContext(ctx, now)
		if err != nil {
			return err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return nil
		}
	}
	return nil
}
