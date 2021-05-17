package pg

import (
	"bytes"
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
				writeError(w, err)
				return
			}
			tx, err := pg.db.BeginTxx(ctx, nil)
			if err != nil {
				log.Error().Err(err).Msg("starting pg transaction failed")
				writeError(w, err)
				return
			}
			ctx = setTx(ctx, tx)

			// this response wrapper hold onto the response until after the
			// handler has executed and we know its status. if the commit
			// fails, we want an error response.
			fw := newResponseWriter(w)
			ww := middleware.NewWrapResponseWriter(fw, r.ProtoMajor)
			next.ServeHTTP(ww, r.WithContext(ctx))

			failed := false
			rescode := ww.Status()
			if statusFailed(rescode) {
				failed = true
				log.Debug().Msg("rollback")
				err = tx.Rollback()
			} else {
				log.Debug().Msg("commit")
				err = tx.Commit()
			}

			if !failed && err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			if err != nil {
				log.Error().Err(err).Msg("commit/rollback failed")
			}
			if err := fw.flush(); err != nil {
				log.Error().Err(err).Msg("flush failed")
			}
		}

		return http.HandlerFunc(fn)
	}
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
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

// responseWriter holds off on sending a response until flush is called.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	out        *bytes.Buffer
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		out:            &bytes.Buffer{},
	}
}

func (w *responseWriter) WriteHeader(statusCode int) {
	if w.statusCode == 0 {
		w.statusCode = statusCode
	}
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return w.out.Write(b)
}

func (w *responseWriter) flush() error {
	code := w.statusCode
	if code == 0 {
		code = 200
	}
	w.ResponseWriter.WriteHeader(code)
	_, err := w.ResponseWriter.Write(w.out.Bytes())
	return err
}
