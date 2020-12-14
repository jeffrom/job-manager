package bepostgres

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
)

type contextKey string

const txKey = contextKey("tx")

// Middleware provides transaction middleware.
func (pg *Postgres) Middleware() func(next http.Handler) http.Handler {
	log := pg.cfg.Logger
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if err := pg.ensureConn(ctx); err != nil {
				log.Error().Err(err).Msg("getting pg conn failed")
				return
			}
			tx, err := pg.db.BeginTxx(ctx, nil)
			if err != nil {
				log.Error().Err(err).Msg("starting pg transaction failed")
				return
			}
			ctx = setTx(ctx, tx)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r.WithContext(ctx))

			if statusFailed(ww.Status()) {
				log.Debug().Msg("rollback")
				err = tx.Rollback()
			} else {
				log.Debug().Msg("commit")
				err = tx.Commit()
			}

			if err != nil {
				log.Error().Err(err).Msg("commit/rollback failed")
				return
			}
		}

		return http.HandlerFunc(fn)
	}
}

func statusFailed(status int) bool {
	return status != 0 && status < 200 || status >= 300
}

func setTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func getTx(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(txKey).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}
