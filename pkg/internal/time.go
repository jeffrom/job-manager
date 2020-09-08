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
		C: make(chan time.Time, 100),
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
	t.C <- now
}
