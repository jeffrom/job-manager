// Package job contains the protocol for working with jobs and queues.
package v1

import (
	uuid "github.com/satori/go.uuid"

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

func FromProto(msg *Job) *resource.Job {
	return &resource.Job{
		ID: msg.Id,
		// Version: resource.NewVersion(msg.V),
	}
}
