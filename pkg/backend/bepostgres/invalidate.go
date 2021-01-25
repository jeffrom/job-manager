package bepostgres

import (
	"context"
)

func (pg *Postgres) InvalidateJobs(ctx context.Context) error {
	if err := pg.ensureConn(ctx); err != nil {
		return err
	}

	return nil
}
