package v1

import (
	"database/sql"
	"time"

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
	deletedAt := sql.NullTime{Valid: isDeleted, Time: msg.DeletedAt.AsTime()}

	return &resource.Queue{
		Name:            msg.Id,
		Version:         resource.NewVersion(msg.V),
		Retries:         int(msg.Retries),
		Duration:        resource.Duration(msg.Duration.AsDuration()),
		ClaimDuration:   resource.Duration(msg.ClaimDuration.AsDuration()),
		CheckinDuration: resource.Duration(msg.CheckinDuration.AsDuration()),
		Unique:          msg.Unique,
		Labels:          label.Labels(msg.Labels),
		SchemaRaw:       msg.Schema,
		BackoffInitial:  resource.Duration(msg.BackoffInitialDuration.AsDuration()),
		BackoffMax:      resource.Duration(msg.BackoffMaxDuration.AsDuration()),
		BackoffFactor:   msg.BackoffFactor,
		CreatedAt:       msg.CreatedAt.AsTime(),
		UpdatedAt:       msg.UpdatedAt.AsTime(),
		DeletedAt:       deletedAt,
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
		Id:                     res.Name,
		V:                      res.Version.Raw(),
		Retries:                int32(res.Retries),
		Duration:               durationpb.New(time.Duration(res.Duration)),
		ClaimDuration:          durationpb.New(time.Duration(res.ClaimDuration)),
		CheckinDuration:        durationpb.New(time.Duration(res.CheckinDuration)),
		Unique:                 res.Unique,
		Labels:                 res.Labels,
		Schema:                 res.SchemaRaw,
		BackoffInitialDuration: durationpb.New(time.Duration(res.BackoffInitial)),
		BackoffMaxDuration:     durationpb.New(time.Duration(res.BackoffMax)),
		BackoffFactor:          res.BackoffFactor,
		CreatedAt:              timestamppb.New(res.CreatedAt),
		UpdatedAt:              timestamppb.New(res.UpdatedAt),
		DeletedAt:              deletedAt,
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
