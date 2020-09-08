// Package v1 contains the v1 client API.
package v1

import (
	"context"
)

type Interface interface {
	Ping(ctx context.Context) error

	EnqueueJobs(ctx context.Context, msg *apiv1.EnqueueRequest) (*apiv1.EnqueueResponse, error)
	DequeueJobs(ctx context.Context, msg *apiv1.DequeueRequest) (*apiv1.DequeueResponse, error)
	AckJobs(ctx context.Context, msg *apiv1.AckRequest) (*apiv1.AckResponse, error)

	SaveQueue(ctx context.Context, msg *apiv1.SaveQueueParams) (*apiv1.SaveQueueResponse, error)
	ListQueues(ctx context.Context, msg *apiv1.ListQueuesRequest) (*apiv1.ListQueuesResponse, error)
	GetQueue(ctx context.Context, msg *apiv1.GetQueueRequest) (*apiv1.GetQueueResponse, error)
}
