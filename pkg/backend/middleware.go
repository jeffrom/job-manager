package backend

import (
	"context"
	"net/http"
)

type contextKey string

var backendContextKey contextKey = "backend"

func Middleware(be Interface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), backendContextKey, be)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func FromMiddleware(ctx context.Context) Interface {
	return ctx.Value(backendContextKey).(Interface)
}
