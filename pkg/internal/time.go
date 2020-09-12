package internal

import (
	"context"
	"time"
)

type contextKey string

var timeContextKey contextKey = "time"
var tickerContextKey contextKey = "time.ticker"

func SetTimeProvider(ctx context.Context, t TimeProvider) context.Context {
	return context.WithValue(ctx, timeContextKey, t)
}

func SetMockTime(ctx context.Context, nows ...time.Time) context.Context {
	return SetTimeProvider(ctx, &MockTime{nows: nows})
}

func SetTicker(ctx context.Context, tick Ticker) context.Context {
	return context.WithValue(ctx, tickerContextKey, tick)
}

func GetTimeProvider(ctx context.Context) TimeProvider {
	p, ok := ctx.Value(timeContextKey).(TimeProvider)
	if ok {
		return p
	}
	return defaultTimeProvider
}

func GetTicker(ctx context.Context) Ticker {
	return ctx.Value(tickerContextKey).(Ticker)
}

type TimeProvider interface {
	Now() time.Time
}

var defaultTimeProvider = Time(0)

type Time int

func (t Time) Now() time.Time { return time.Now() }

type Ticker interface {
	Chan() <-chan time.Time
	Stop()
}

type defaultTicker struct {
	*time.Ticker
}

func (t *defaultTicker) Chan() <-chan time.Time { return t.Ticker.C }

func NewTicker(d time.Duration) *defaultTicker {
	return &defaultTicker{
		Ticker: time.NewTicker(d),
	}
}

type MockTime struct {
	nows []time.Time
}

func (t *MockTime) SetNow(nows ...time.Time) { t.nows = nows }
func (t *MockTime) Now() time.Time {
	if len(t.nows) == 0 {
		panic("no more times stored")
	}
	now := t.nows[0]
	if len(t.nows) > 1 {
		t.nows = t.nows[1:]
	}
	return now
}

type MockTick struct {
	nows []time.Time
	C    chan time.Time
}

func NewMockTick(d time.Duration) *MockTick {
	return &MockTick{
		C: make(chan time.Time),
	}
}

func (t *MockTick) Chan() <-chan time.Time   { return t.C }
func (t *MockTick) SetNow(nows ...time.Time) { t.nows = nows }
func (t *MockTick) Stop()                    { close(t.C) }
func (t *MockTick) Tick() {
	if len(t.nows) == 0 {
		panic("no more times stored")
	}
	now := t.nows[0]
	if len(t.nows) > 1 {
		t.nows = t.nows[1:]
	}

	// i believe the stdlib ticker also throws away any pending message before
	// firing the tick.
	select {
	case <-t.C:
	default:
	}
	t.C <- now
}
