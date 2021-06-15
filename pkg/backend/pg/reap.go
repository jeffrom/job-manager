package pg

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/jeffrom/job-manager/mjob/resource"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/internal"
)

// Reap cleans up old jobs and queue versions.
func (pg *Postgres) Reap(ctx context.Context, cfg *backend.ReaperConfig) error {
	if err := pg.ensureConn(ctx); err != nil {
		return err
	}

	for {
		if done, err := pg.reapOne(ctx); err != nil {
			return err
		} else if done {
			return nil
		}
		internal.IgnoreError(internal.Sleep(ctx, 2*time.Second))
	}
}

func (pg *Postgres) reapOne(ctx context.Context) (bool, error) {
	now := internal.GetTimeProvider(ctx).Now().UTC()
	tx, err := pg.db.BeginTxx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()
	q := "SELECT * FROM jobs WHERE enqueued_at < $1 ORDER BY id ASC LIMIT 50 FOR UPDATE SKIP LOCKED"
	var rows []*resource.Job
	err = sqlx.SelectContext(ctx, tx, &rows, q, now.Add(-pg.cfg.ReapAge))
	if err != nil {
		return false, err
	}
	if len(rows) == 0 {
		return true, nil
	}

	for _, jb := range rows {
		fmt.Println("reap", jb.ID)
		if err := pg.deleteJobUniqueness(ctx, tx, jb); err != nil {
			return false, err
		}
		if err := pg.deleteJobResults(ctx, tx, jb); err != nil {
			return false, err
		}
		if err := pg.deleteJobClaims(ctx, tx, jb); err != nil {
			return false, err
		}
		if err := pg.deleteJobCheckins(ctx, tx, jb); err != nil {
			return false, err
		}
	}
	if err := pg.deleteJobs(ctx, tx, rows); err != nil {
		return false, err
	}

	return false, tx.Commit()
}

func (pg *Postgres) deleteJobUniqueness(ctx context.Context, tx *sqlx.Tx, jb *resource.Job) error {
	key, err := jb.ArgKey()
	if err != nil {
		return err
	}
	q := "DELETE FROM job_uniqueness WHERE key = $1"
	_, err = tx.ExecContext(ctx, q, []byte(key))
	return err
}

func (pg *Postgres) deleteJobResults(ctx context.Context, tx *sqlx.Tx, jb *resource.Job) error {
	q := "DELETE FROM job_results WHERE job_id = $1"
	_, err := tx.ExecContext(ctx, q, jb.ID)
	return err
}

func (pg *Postgres) deleteJobClaims(ctx context.Context, tx *sqlx.Tx, jb *resource.Job) error {
	q := "DELETE FROM job_claims WHERE job_id = $1"
	_, err := tx.ExecContext(ctx, q, jb.ID)
	return err
}

func (pg *Postgres) deleteJobCheckins(ctx context.Context, tx *sqlx.Tx, jb *resource.Job) error {
	q := "DELETE FROM job_checkins WHERE job_id = $1"
	_, err := tx.ExecContext(ctx, q, jb.ID)
	return err
}

func (pg *Postgres) deleteJobs(ctx context.Context, tx *sqlx.Tx, jobs []*resource.Job) error {
	var args []interface{}
	ids := make([]int64, len(jobs))
	for i, jb := range jobs {
		n, err := strconv.ParseInt(jb.ID, 10, 64)
		if err != nil {
			return err
		}
		ids[i] = n
	}
	args = append(args, ids)
	q := "DELETE FROM jobs WHERE id IN (?)"

	var err error
	q, args, err = sqlx.In(q, args...)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, tx.Rebind(q), args...)
	return err
}
