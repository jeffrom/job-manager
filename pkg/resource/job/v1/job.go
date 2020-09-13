// Package v1 contains the protocol for working with jobs and queues.
package v1

import (
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jeffrom/job-manager/pkg/label"
	"github.com/jeffrom/job-manager/pkg/resource"
)

const (
	StatusUnknown   = Status_STATUS_UNSPECIFIED
	StatusQueued    = Status_STATUS_QUEUED
	StatusRunning   = Status_STATUS_RUNNING
	StatusComplete  = Status_STATUS_COMPLETE
	StatusDead      = Status_STATUS_DEAD
	StatusCancelled = Status_STATUS_CANCELLED
	StatusInvalid   = Status_STATUS_INVALID
	StatusFailed    = Status_STATUS_FAILED
)

func NewID() string {
	return uuid.NewV4().String()
}

func HasStatus(job *Job, statuses []Status) bool {
	for _, st := range statuses {
		if job.Status == st {
			return true
		}
	}
	return false
}

func IsComplete(status Status) bool {
	switch status {
	case StatusComplete, StatusCancelled, StatusInvalid, StatusDead:
		return true
	}
	return false
}

func NewJobFromResource(jb *resource.Job) (*Job, error) {
	v, err := structpb.NewList(jb.Args)
	if err != nil {
		return nil, err
	}
	var data *Data
	if jb.Data != nil {
		data := &Data{}
		if jb.Data.Data != nil {
			// fmt.Printf("%#v\n", jb.Data.Data)
			datav, err := structpb.NewValue(jb.Data.Data)
			if err != nil {
				return nil, err
			}
			data.Data = datav
		}
		if len(jb.Data.Claims) > 0 {
			data.Claims = jb.Data.Claims.Format()
		}
	}

	results, err := jobResultsToProto(jb.Results)
	if err != nil {
		return nil, err
	}
	return &Job{
		Id:         jb.ID,
		V:          jb.Version.Raw(),
		Name:       jb.Name,
		QueueV:     jb.QueueVersion.Raw(),
		Args:       v.Values,
		Data:       data,
		Status:     Status(jb.Status),
		Attempt:    int32(jb.Attempt),
		Checkins:   jobCheckinsToProto(jb.Checkins),
		Results:    results,
		EnqueuedAt: timestamppb.New(jb.EnqueuedAt),
	}, nil
}

func NewJobsFromResources(jbs []*resource.Job) ([]*Job, error) {
	jobs := make([]*Job, len(jbs))
	for i, j := range jbs {
		pj, err := NewJobFromResource(j)
		if err != nil {
			return nil, err
		}
		jobs[i] = pj
	}
	return jobs, nil
}

func NewJobFromProto(msg *Job, claims label.Claims) *resource.Job {
	lv := &structpb.ListValue{Values: msg.Args}
	var data *resource.JobData
	if msg.Data != nil || len(claims) > 0 {
		data = &resource.JobData{}
		if msg.Data.Data != nil {
			data.Data = msg.Data.Data.AsInterface()
		}
		if len(claims) > 0 {
			data.Claims = claims
		}
	}

	return &resource.Job{
		ID:           msg.Id,
		Version:      resource.NewVersion(msg.V),
		Name:         msg.Name,
		QueueVersion: resource.NewVersion(msg.QueueV),
		Args:         lv.AsSlice(),
		Data:         data,
		Status:       jobStatusFromProto(msg.Status),
		Attempt:      int(msg.Attempt),
		Checkins:     jobCheckinsFromProto(msg.Checkins),
		Results:      jobResultsFromProto(msg.Results),
		EnqueuedAt:   msg.EnqueuedAt.AsTime(),
	}
}

func NewJobsFromProto(msgs []*Job) ([]*resource.Job, error) {
	jobs := make([]*resource.Job, len(msgs))
	for i, msg := range msgs {
		jb := NewJobFromProto(msg, nil)
		if msg.Data != nil && len(msg.Data.Claims) > 0 {
			claims, err := label.ParseClaims(msg.Data.Claims)
			if err != nil {
				return nil, err
			}
			if jb.Data == nil {
				jb.Data = &resource.JobData{}
			}
			jb.Data.Claims = claims
		}
		jobs[i] = jb
	}

	return jobs, nil
}

func jobStatusFromProto(status Status) resource.Status {
	switch status {
	case Status_STATUS_UNSPECIFIED:
		return resource.StatusUnspecified
	case Status_STATUS_QUEUED:
		return resource.StatusQueued
	case Status_STATUS_RUNNING:
		return resource.StatusRunning
	case Status_STATUS_COMPLETE:
		return resource.StatusComplete
	case Status_STATUS_DEAD:
		return resource.StatusDead
	case Status_STATUS_CANCELLED:
		return resource.StatusCancelled
	case Status_STATUS_INVALID:
		return resource.StatusInvalid
	case Status_STATUS_FAILED:
		return resource.StatusFailed
	default:
		panic("job/v1: unknown status")
	}
}

func jobCheckinsFromProto(checkins []*Checkin) []*resource.JobCheckin {
	rcs := make([]*resource.JobCheckin, len(checkins))
	for i, c := range checkins {
		rcs[i] = &resource.JobCheckin{
			Data:      c.Data.AsInterface(),
			CreatedAt: c.CreatedAt.AsTime(),
		}
	}
	return rcs
}

func jobCheckinsToProto(pcheckins []*resource.JobCheckin) []*Checkin {
	cs := make([]*Checkin, len(pcheckins))
	for i, pc := range pcheckins {
		v, err := structpb.NewValue(pc.Data)
		if err != nil {
			return nil
		}
		cs[i] = &Checkin{
			Data:      v,
			CreatedAt: timestamppb.New(pc.CreatedAt),
		}
	}
	return cs
}

func jobResultsFromProto(results []*Result) []*resource.JobResult {
	jrs := make([]*resource.JobResult, len(results))
	for i, r := range results {
		jrs[i] = &resource.JobResult{
			Attempt:     int(r.Attempt),
			Status:      jobStatusFromProto(r.Status),
			Data:        r.Data.AsInterface(),
			StartedAt:   r.StartedAt.AsTime(),
			CompletedAt: r.CompletedAt.AsTime(),
		}
	}
	return jrs
}

func jobResultsToProto(results []*resource.JobResult) ([]*Result, error) {
	prs := make([]*Result, len(results))
	for i, pr := range results {
		v, err := structpb.NewValue(pr.Data)
		if err != nil {
			return nil, err
		}
		prs[i] = &Result{
			Attempt:     int32(pr.Attempt),
			Status:      Status(pr.Status),
			Data:        v,
			StartedAt:   timestamppb.New(pr.StartedAt),
			CompletedAt: timestamppb.New(pr.CompletedAt),
		}
	}
	return prs, nil
}

func MarshalJob(jb *resource.Job) ([]byte, error) {
	jbp, err := NewJobFromResource(jb)
	if err != nil {
		return nil, err
	}
	b, err := proto.Marshal(jbp)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UnmarshalJob(b []byte, qmsg *Job) (*resource.Job, error) {
	if qmsg == nil {
		qmsg = &Job{}
	}
	if err := proto.Unmarshal(b, qmsg); err != nil {
		return nil, err
	}
	var claims label.Claims
	if qmsg.Data != nil && len(qmsg.Data.Claims) > 0 {
		var err error
		claims, err = label.ParseClaims(qmsg.Data.Claims)
		if err != nil {
			return nil, err
		}
	}
	return NewJobFromProto(qmsg, claims), nil
}
