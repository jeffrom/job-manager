package backend

import (
	"context"
	"time"
)

// Reaper adds functionality to clean up old job data.
type Reaper interface {
	Reap(ctx context.Context, cfg *ReaperConfig) error
}

type ReaperConfig struct {
	Interval    time.Duration `json:"reap_interval"`
	Age         time.Duration `json:"reap_age"`
	Max         int           `json:"reap_max"`
	MaxVersions int           `json:"reap_max_versions"`
}

type NoopReaper struct{}

func (r *NoopReaper) Reap(ctx context.Context, cfg *ReaperConfig) error {
	return nil
}
