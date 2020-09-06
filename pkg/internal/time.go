package internal

import (
	"time"
)

type TimeProvider interface {
	Now() time.Time
}

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
	t time.Time
}

func (t *MockTime) Now() time.Time      { return t.t }
func (t *MockTime) SetNow(ti time.Time) { t.t = ti }

type MockTick struct {
	C chan time.Time
}

func NewMockTick(d time.Duration) *MockTick {
	return &MockTick{
		C: make(chan time.Time),
	}
}

func (t *MockTick) Chan() <-chan time.Time { return t.C }
func (t *MockTick) Tick(ti time.Time)      { t.C <- ti }
func (t *MockTick) Stop()                  { close(t.C) }
