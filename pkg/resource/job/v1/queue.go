package v1

import (
	"github.com/jeffrom/job-manager/pkg/label"
	"github.com/jeffrom/job-manager/pkg/resource"
)

func NewQueueFromProto(msg *Queue) *resource.Queue {
	return &resource.Queue{
		ID:              msg.Id,
		Version:         resource.NewVersion(msg.V),
		Concurrency:     int(msg.Concurrency),
		Retries:         int(msg.Retries),
		Duration:        msg.Duration.AsDuration(),
		CheckinDuration: msg.CheckinDuration.AsDuration(),
		Unique:          msg.Unique,
		Labels:          label.Labels(msg.Labels),
		SchemaRaw:       msg.Schema,
		CreatedAt:       msg.CreatedAt.AsTime(),
		UpdatedAt:       msg.UpdatedAt.AsTime(),
		DeletedAt:       msg.DeletedAt.AsTime(),
	}
}
