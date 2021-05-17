package pg

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/jeffrom/job-manager/pkg/internal"
)

func (pg *Postgres) InvalidateJobs(ctx context.Context) error {
	if err := pg.ensureConn(ctx); err != nil {
		return err
	}

	c := pg.db
	sql := "UPDATE jobs SET status = 'failed', v=v+1, completed_at = ? FROM (SELECT jobs.id, queues.retries FROM jobs LEFT JOIN queues ON jobs.queue_id = queues.id WHERE status = 'running' AND jobs.started_at + (queues.duration / 1000) * interval '1 microsecond' < ? AND jobs.attempt <= queues.retries ORDER BY jobs.id ASC LIMIT 50 FOR UPDATE OF jobs SKIP LOCKED) AS to_update WHERE jobs.id = to_update.id"
	q, err := sqlx.PreparexContext(ctx, c, c.Rebind(sql))
	if err != nil {
		return err
	}
	defer q.Close()

	now := internal.GetTimeProvider(ctx).Now().UTC()
	if err := execJobUpdateLoop(ctx, q, now); err != nil {
		return err
	}

	sql = "UPDATE jobs SET status = 'dead', v=v+1, completed_at = ? FROM (SELECT jobs.id, queues.retries FROM jobs LEFT JOIN queues ON jobs.queue_id = queues.id WHERE status = 'running' AND jobs.started_at + (queues.duration / 1000) * interval '1 microsecond' < ? AND jobs.attempt > queues.retries ORDER BY jobs.id ASC LIMIT 50 FOR UPDATE OF jobs SKIP LOCKED) AS to_update WHERE jobs.id = to_update.id"
	q, err = sqlx.PreparexContext(ctx, c, c.Rebind(sql))
	if err != nil {
		return err
	}
	defer q.Close()

	if err := execJobUpdateLoop(ctx, q, now); err != nil {
		return err
	}
	return nil
}

func execJobUpdateLoop(ctx context.Context, stmt *sqlx.Stmt, now time.Time) error {
	for {
		res, err := stmt.ExecContext(ctx, now, now)
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
