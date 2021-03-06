package internal

import (
	"context"
	"time"
)

var MockTimeKey ContextKey = "mocktime"

func GetMockTime(ctx context.Context) *time.Time {
	t := ctx.Value(MockTimeKey)
	if t == nil {
		return nil
	}
	return t.(*time.Time)
}
