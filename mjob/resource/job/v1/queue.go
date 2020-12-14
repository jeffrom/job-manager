package v1

import (
	"database/sql"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jeffrom/job-manager/mjob/label"
	"github.com/jeffrom/job-manager/mjob/resource"
)

func NewQueueFromProto(msg *Queue) *resource.Queue {
	if msg == nil {
		return nil
	}
	isDeleted := msg.DeletedAt != nil
	return &resource.Queue{
		Name:            msg.Id,
		Version:         resource.NewVersion(msg.V),
		Concurrency:     int(msg.Concurrency),
		Retries:         int(msg.Retries),
		Duration:        msg.Duration.AsDuration(),
		ClaimDuration:   msg.ClaimDuration.AsDuration(),
		CheckinDuration: msg.CheckinDuration.AsDuration(),
		Unique:          msg.Unique,
		Labels:          label.Labels(msg.Labels),
		SchemaRaw:       msg.Schema,
		CreatedAt:       msg.CreatedAt.AsTime(),
		UpdatedAt:       msg.UpdatedAt.AsTime(),
		DeletedAt:       sql.NullTime{Valid: isDeleted, Time: msg.DeletedAt.AsTime()},
	}
}

func NewQueuesFromProto(msgs []*Queue) []*resource.Queue {
	qs := make([]*resource.Queue, len(msgs))
	for i, msg := range msgs {
		qs[i] = NewQueueFromProto(msg)
	}
	return qs
}

func NewQueueFromResource(res *resource.Queue) *Queue {
	if res == nil {
		return nil
	}
	var deletedAt *timestamppb.Timestamp
	if res.DeletedAt.Valid {
		deletedAt = timestamppb.New(res.DeletedAt.Time)
	}
	return &Queue{
		Id:              res.Name,
		V:               res.Version.Raw(),
		Concurrency:     int32(res.Concurrency),
		Retries:         int32(res.Retries),
		Duration:        durationpb.New(res.Duration),
		ClaimDuration:   durationpb.New(res.ClaimDuration),
		CheckinDuration: durationpb.New(res.CheckinDuration),
		Unique:          res.Unique,
		Labels:          res.Labels,
		Schema:          res.SchemaRaw,
		CreatedAt:       timestamppb.New(res.CreatedAt),
		UpdatedAt:       timestamppb.New(res.UpdatedAt),
		DeletedAt:       deletedAt,
	}
}

func NewQueuesFromResources(resources []*resource.Queue) []*Queue {
	qs := make([]*Queue, len(resources))
	for i, rq := range resources {
		qs[i] = NewQueueFromResource(rq)
	}
	return qs
}

func MarshalQueue(q *resource.Queue) ([]byte, error) {
	return proto.Marshal(NewQueueFromResource(q))
}

func UnmarshalQueue(b []byte, qmsg *Queue) (*resource.Queue, error) {
	if qmsg == nil {
		qmsg = &Queue{}
	}
	if err := proto.Unmarshal(b, qmsg); err != nil {
		return nil, err
	}
	return NewQueueFromProto(qmsg), nil
}
