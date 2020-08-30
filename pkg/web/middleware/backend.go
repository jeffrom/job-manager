package middleware

import (
	"context"
	"net/http"

	"github.com/jeffrom/job-manager/pkg/backend"
)

var backendContextKey contextKey = "backend"

func Backend(be backend.Interface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), backendContextKey, be)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func GetBackend(ctx context.Context) backend.Interface {
	return ctx.Value(backendContextKey).(backend.Interface)
}
