package backend

import "context"

type Invalidator interface {
	InvalidateJobs(ctx context.Context) error
}

type NoopInvalidator struct{}

func (iv *NoopInvalidator) InvalidateJobs(ctx context.Context) error {
	return nil
}
