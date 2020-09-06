package v1

import (
	"github.com/jeffrom/job-manager/pkg/label"
	"github.com/jeffrom/job-manager/pkg/resource"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func NewQueueFromResource(res *resource.Queue) *Queue {
	return &Queue{
		Id:              res.ID,
		V:               res.Version.Raw(),
		Concurrency:     int32(res.Concurrency),
		Retries:         int32(res.Retries),
		Duration:        durationpb.New(res.Duration),
		CheckinDuration: durationpb.New(res.CheckinDuration),
		Unique:          res.Unique,
		Labels:          res.Labels,
		Schema:          res.SchemaRaw,
		CreatedAt:       timestamppb.New(res.CreatedAt),
		UpdatedAt:       timestamppb.New(res.UpdatedAt),
		DeletedAt:       timestamppb.New(res.DeletedAt),
	}
}

func NewQueuesFromResources(resources []*resource.Queue) []*Queue {
	qs := make([]*Queue, len(resources))
	for i, rq := range resources {
		qs[i] = NewQueueFromResource(rq)
	}
	return qs
}
