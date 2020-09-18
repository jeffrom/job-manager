package bepg

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
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if err := pg.ensureConn(ctx); err != nil {
				panic(err)
			}
			tx, err := pg.db.BeginTxx(ctx, nil)
			if err != nil {
				panic(err)
				// return
			}
			ctx = setTx(ctx, tx)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r.WithContext(ctx))

			if statusFailed(ww.Status()) {
				err = tx.Rollback()
			} else {
				err = tx.Commit()
			}

			if err != nil {
				panic(err)
			}
		}

		return http.HandlerFunc(fn)
	}
}

func statusFailed(status int) bool {
	return status < 200 || status >= 300
}

func setTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func getTx(ctx context.Context) *sqlx.Tx {
	return ctx.Value(txKey).(*sqlx.Tx)
}
