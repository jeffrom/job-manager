package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jeffrom/job-manager/pkg/internal"
)

const (
	mockTimeHeader = "fake-time"
)

func Time(t internal.TimeProvider, tick internal.Ticker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// fmt.Println("mw mockTimeHeader:", r.Header.Get(mockTimeHeader))
			ctx := r.Context()
			if h := r.Header.Get(mockTimeHeader); h != "" {
				times := strings.Split(h, ",")
				allts := make([]time.Time, len(times))
				for i, hts := range times {
					uts, err := strconv.ParseInt(hts, 10, 64)
					if err != nil {
						panic(err)
					}
					ts := time.Unix(uts, 0).UTC()
					// fmt.Println("mw:", ts.Format(time.Stamp))
					allts[i] = ts
				}

				mt := &internal.MockTime{}
				mt.SetNow(allts...)
				t = mt

				mtick := internal.NewMockTick(0)
				mtick.SetNow(allts...)
				defer mtick.Stop()
				tick = mtick
			}
			ctx = internal.SetTimeProvider(ctx, t)
			ctx = internal.SetTicker(ctx, tick)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
