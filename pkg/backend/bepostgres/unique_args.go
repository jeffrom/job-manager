package bepostgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jeffrom/job-manager/pkg/internal"
	"github.com/jmoiron/sqlx"
)

func (pg *Postgres) GetJobUniqueArgs(ctx context.Context, keys []string) ([]string, bool, error) {
	c, err := pg.getConn(ctx)
	if err != nil {
		return nil, false, err
	}

	q := "SELECT job_id FROM job_uniqueness WHERE key IN (?)"
	args := stringsToBytea(keys)
	q, iargs, err := sqlx.In(q, args)
	if err != nil {
		return nil, false, err
	}
	rows, err := c.QueryxContext(ctx, c.Rebind(q), iargs...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, false, err
	}
	var ids []string
	for rows.Next() {
		id := ""
		if err := rows.Scan(&id); err != nil {
			return nil, true, err
		}
		ids = append(ids, id)
	}
	if len(ids) > 0 {
		return ids, true, nil
	}
	return nil, false, nil
}

func (pg *Postgres) SetJobUniqueArgs(ctx context.Context, ids, keys []string) error {
	if len(ids) != len(keys) {
		panic("backend/postgres: mismatched ids, key args")
	}
	c, err := pg.getConn(ctx)
	if err != nil {
		return err
	}
	stmt, err := c.PrepareContext(ctx, "INSERT INTO job_uniqueness (job_id, key, created_at) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := internal.GetTimeProvider(ctx).Now().UTC()
	for i, arg := range stringsToBytea(keys) {
		if _, err := stmt.ExecContext(ctx, ids[i], arg, now); err != nil {
			return err
		}
	}
	return nil
}

func (pg *Postgres) DeleteJobUniqueArgs(ctx context.Context, ids, keys []string) error {
	c, err := pg.getConn(ctx)
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		stmt, err := c.PrepareContext(ctx, "DELETE FROM job_uniqueness WHERE key = $1")
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, key := range stringsToBytea(keys) {
			if _, err := stmt.ExecContext(ctx, key); err != nil && !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}
	}

	if len(ids) > 0 {
		stmt, err := c.PrepareContext(ctx, "DELETE FROM job_uniqueness WHERE job_id = $1")
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, id := range ids {
			if _, err := stmt.ExecContext(ctx, id); err != nil && !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}
	}
	return nil
}
