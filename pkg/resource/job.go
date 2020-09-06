package resource

import (
	"time"

	"github.com/jeffrom/job-manager/pkg/label"
)

type Job struct {
	ID           string        `json:"id"`
	Version      *Version      `json:"version"`
	Name         string        `json:"name"`
	QueueVersion *Version      `json:"queue_version"`
	Args         []interface{} `json:"args"`
	Data         *JobData      `json:"data,omitempty"`
	Status       Status        `json:"status"`
	Attempt      int           `json:"attempt,omitempty"`
	Checkins     []*JobCheckin `json:"checkins,omitempty"`
	Results      []*JobResult  `json:"results,omitempty"`
	EnqueuedAt   time.Time     `json:"enqueued_at,omitempty"`
}

func (jb *Job) IsAttempted() bool {
	return StatusIsAttempted(jb.Status)
}

func (jb *Job) HasStatus(statuses ...Status) bool {
	for _, st := range statuses {
		if jb.Status == st {
			return true
		}
	}
	return false
}

type Jobs struct {
	Jobs []*Job `json:"jobs"`
}

type JobData struct {
	Data interface{} `json:"data,omitempty"`
}

type JobCheckin struct {
	Data      interface{} `json:"data,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}

type JobResult struct {
	Attempt     int         `json:"attempt"`
	Status      Status      `json:"status"`
	Data        interface{} `json:"data,omitempty"`
	StartedAt   time.Time   `json:"started_at"`
	CompletedAt time.Time   `json:"completed_at"`
}

type JobListParams struct {
	Names         []string        `json:"names,omitempty"`
	Statuses      []Status        `json:"statuses,omitempty"`
	Selectors     label.Selectors `json:"selectors,omitempty"`
	EnqueuedSince time.Time       `json:"enqueued_since,omitempty"`
	EnqueuedUntil time.Time       `json:"enqueued_until,omitempty"`
}
