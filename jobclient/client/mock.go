package client

import (
	"context"
	"time"

	"github.com/jeffrom/job-manager/jobclient/internal"
)

func SetMockTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, internal.MockTimeKey, &t)
}
