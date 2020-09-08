package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jeffrom/job-manager/pkg/internal"
)

var timeContextKey contextKey = "time"
var tickerContextKey contextKey = "time.ticker"

const (
	mockTimeHeader   = "fake-time"
	mockTickerHeader = "fake-ticker"
)

func Time(t internal.TimeProvider, tick internal.Ticker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if h := r.Header.Get(mockTimeHeader); h != "" {
				times := strings.Split(h, ",")
				allTs := make([]time.Time, len(times))
				for i, hts := range times {
					uts, err := strconv.ParseInt(hts, 10, 64)
					if err != nil {
						panic(err)
						return
					}
					allTs[i] = time.Unix(uts, 0)
				}

				mt := &internal.MockTime{}
				mt.SetNow(allTs...)
				t = mt

				mtick := internal.NewMockTick(0)
				mtick.SetNow(allTs...)
				defer mtick.Stop()
				tick = mtick
			}
			ctx = context.WithValue(ctx, timeContextKey, t)
			ctx = context.WithValue(ctx, tickerContextKey, tick)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func SetTime(ctx context.Context, t internal.TimeProvider) context.Context {
	return context.WithValue(ctx, timeContextKey, t)
}

func SetTicker(ctx context.Context, tick internal.Ticker) context.Context {
	return context.WithValue(ctx, tickerContextKey, tick)
}

func GetTime(ctx context.Context) internal.TimeProvider {
	return ctx.Value(timeContextKey).(internal.TimeProvider)
}

func GetTicker(ctx context.Context) internal.Ticker {
	return ctx.Value(tickerContextKey).(internal.Ticker)
}
