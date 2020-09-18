// Package logger contains a logger for use by other job-manager packages.
package logger

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

type contextKey string

var loggerKey = contextKey("logger")
var reqLogKey = contextKey("reqLog")
var queryKey = contextKey("query")

type Logger struct {
	zerolog.Logger
}

func FromContext(ctx context.Context) *Logger {
	return ctx.Value(loggerKey).(*Logger)
}

func RequestLogFromContext(ctx context.Context) *zerolog.Event {
	return ctx.Value(reqLogKey).(*zerolog.Event)
}

func New(l zerolog.Logger) *Logger { return &Logger{Logger: l} }

func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		rec := httptest.NewRecorder()

		// req, _ := httputil.DumpRequest(r, r.Method != http.MethodGet)

		ctx := r.Context()

		reqID := ctx.Value(middleware.RequestIDKey).(string)
		l.Debug().
			Str("req_id", reqID).
			Msg("start request")

		reqLog := l.Info()
		reqLog.Timestamp().
			Str("method", r.Method).
			Str("path", r.URL.EscapedPath()).
			Str("ip", r.RemoteAddr).
			Str("req_id", reqID)
			//.Bytes("request", req)

		// var body string
		var queries []string

		defer func(begin time.Time) {
			status := ww.Status()
			reqLog.
				Int64("took", time.Since(begin).Milliseconds()).
				Int("status", status)

			query := r.URL.Query().Encode()
			if query != "" {
				reqLog.Str("query", query)
			}
			// if len(queries) > 0 {
			// 	reqLog.Strs("queries", queries)
			// }

			// if status != http.StatusNotFound {
			// 	reqLog.Str("response", body)
			// }

			// if status >= 500 && status < 600 {
			reqLog.Msg("request")
			// }

		}(time.Now())
		reqLogger := l.With().Str("req_id", reqID).Logger()
		ctx = context.WithValue(ctx, loggerKey, &Logger{Logger: reqLogger})
		ctx = context.WithValue(ctx, reqLogKey, reqLog)
		ctx = context.WithValue(ctx, queryKey, &queries)

		next.ServeHTTP(rec, r.WithContext(ctx))

		// this copies the recorded response to the response writer
		for k, v := range rec.Header() {
			ww.Header()[k] = v
		}

		// body = rec.Body.String()
		ww.WriteHeader(rec.Code)
		rec.Body.WriteTo(ww)
	})
}