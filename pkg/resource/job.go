package resource

import (
	"encoding/json"
	"time"

	"github.com/jeffrom/job-manager/pkg/label"
)

type Job struct {
	ID           string        `json:"id"`
	QueueID      int64         `json:"-" db:"queue"`
	Version      *Version      `json:"version" db:"v"`
	Name         string        `json:"name"`
	QueueVersion *Version      `json:"queue_version" db:"queue_v"`
	Args         []interface{} `json:"args"`
	Data         *JobData      `json:"data,omitempty"`
	Status       *Status       `json:"status"`
	Attempt      int           `json:"attempt,omitempty"`
	Checkins     []*JobCheckin `json:"checkins,omitempty"`
	Results      []*JobResult  `json:"results,omitempty"`
	EnqueuedAt   time.Time     `json:"enqueued_at,omitempty" db:"enqueued_at"`
}

func (jb *Job) String() string {
	b, _ := json.Marshal(jb)
	return string(b)
}

func (jb *Job) Copy() *Job {
	cp := &Job{}
	*cp = *jb
	if jb.Version != nil {
		cp.Version = NewVersion(jb.Version.Raw())
	}
	if jb.QueueVersion != nil {
		cp.QueueVersion = NewVersion(jb.QueueVersion.Raw())
	}
	if jb.Data != nil {
		cp.Data = &JobData{}
		*cp.Data = *jb.Data
	}
	if jb.Checkins != nil {
		cp.Checkins = make([]*JobCheckin, len(jb.Checkins))
		copy(cp.Checkins, jb.Checkins)
	}
	if jb.Results != nil {
		cp.Results = make([]*JobResult, len(jb.Results))
		copy(cp.Results, jb.Results)
	}
	return cp
}

func (jb *Job) LastResult() *JobResult {
	res := jb.Results
	if len(res) == 0 {
		return nil
	}
	return res[len(res)-1]
}

func (jb *Job) LastClaimWindow() time.Time {
	if results := jb.Results; len(results) > 0 {
		for i := len(results) - 1; i >= 0; i-- {
			res := results[i]
			completedAt := res.CompletedAt
			if !completedAt.IsZero() {
				return completedAt
			}
		}
	}
	return jb.EnqueuedAt
}

func (jb *Job) IsAttempted() bool {
	return StatusIsAttempted(jb.Status)
}

func (jb *Job) HasStatus(statuses ...Status) bool {
	for _, st := range statuses {
		if *jb.Status == st {
			return true
		}
	}
	return false
}

type Jobs struct {
	Jobs []*Job `json:"jobs"`
}

func (jobs *Jobs) IDs() []string {
	ids := make([]string, len(jobs.Jobs))
	for i, jb := range jobs.Jobs {
		ids[i] = jb.ID
	}
	return ids
}

func (jobs *Jobs) Queues() []string {
	names := make([]string, len(jobs.Jobs))
	for i, jb := range jobs.Jobs {
		names[i] = jb.Name
	}
	return names
}

type JobData struct {
	Claims label.Claims `json:"claims,omitempty"`
	Data   interface{}  `json:"data,omitempty"`
}

type JobCheckin struct {
	Data      interface{} `json:"data,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}

type JobResult struct {
	Attempt     int         `json:"attempt"`
	Status      *Status     `json:"status"`
	Data        interface{} `json:"data,omitempty"`
	Error       string      `json:"error,omitempty"`
	StartedAt   time.Time   `json:"started_at"`
	CompletedAt time.Time   `json:"completed_at"`
}

type JobListParams struct {
	Names         []string        `json:"names,omitempty"`
	Statuses      []Status        `json:"statuses,omitempty"`
	Selectors     label.Selectors `json:"selectors,omitempty"`
	Claims        label.Claims    `json:"claims,omitempty"`
	EnqueuedSince time.Time       `json:"enqueued_since,omitempty"`
	EnqueuedUntil time.Time       `json:"enqueued_until,omitempty"`
}
