// Package job contains the protocol for working with jobs and queues.
package v1

import (
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/structpb"

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

func NewJobFromProto(msg *Job) *resource.Job {
	lv := &structpb.ListValue{Values: msg.Args}
	return &resource.Job{
		ID:           msg.Id,
		Version:      resource.NewVersion(msg.V),
		Name:         msg.Name,
		QueueVersion: resource.NewVersion(msg.QueueV),
		Args:         lv.AsSlice(),
		Data: &resource.JobData{
			Data: msg.Data.Data.AsInterface(),
		},
		Status:   jobStatusFromProto(msg.Status),
		Attempt:  int(msg.Attempt),
		Checkins: jobCheckinsFromProto(msg.Checkins),
		Results:  jobResultsFromProto(msg.Results),
	}
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
