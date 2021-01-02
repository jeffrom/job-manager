package resource

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jeffrom/job-manager/mjob/label"
)

type Job struct {
	ID           string        `json:"id"`
	QueueID      int64         `json:"-" db:"queue_id"`
	Version      *Version      `json:"version" db:"v"`
	Name         string        `json:"name"`
	QueueVersion *Version      `json:"queue_version" db:"queue_v"`
	Args         []interface{} `json:"args"`
	ArgsRaw      []byte        `json:"-" db:"args"`
	Data         *JobData      `json:"data,omitempty"`
	DataRaw      []byte        `json:"-" db:"data"`
	Status       *Status       `json:"status"`
	Attempt      int           `json:"attempt,omitempty"`
	Checkins     []*JobCheckin `json:"checkins,omitempty"`
	Results      []*JobResult  `json:"results,omitempty"`
	EnqueuedAt   time.Time     `json:"enqueued_at,omitempty" db:"enqueued_at"`
	StartedAt    sql.NullTime  `json:"started_at,omitempty" db:"started_at"`
}

func (jb *Job) String() string {
	b, _ := json.Marshal(jb)
	return string(b)
}

func (jb *Job) Populate() error {
	if jb.Args != nil {
		b, err := json.Marshal(jb.Args)
		if err != nil {
			return err
		}
		jb.ArgsRaw = b
	} else if len(jb.ArgsRaw) > 0 {
		args := []interface{}{}
		if err := json.Unmarshal(jb.ArgsRaw, &args); err != nil {
			return err
		}
		jb.Args = args
	}

	if jb.Data != nil && jb.Data.Data != nil {
		b, err := json.Marshal(jb.Data)
		if err != nil {
			return err
		}
		jb.DataRaw = b
	} else if len(jb.DataRaw) > 0 {
		var data interface{}
		if err := json.Unmarshal(jb.DataRaw, &data); err != nil {
			return err
		}
		if jb.Data == nil {
			jb.Data = &JobData{}
		}
		jb.Data.Data = data
	}
	return nil
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

func (jb *Job) HasStatus(statuses ...*Status) bool {
	for _, st := range statuses {
		if *jb.Status == *st {
			return true
		}
	}
	return false
}

type Jobs struct {
	Jobs []*Job `json:"jobs"`
}

// func (jobs *Jobs) Populate() error {
// 	for _, jb := range jobs {
// 		if err := jb.Populate(); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func (jobs *Jobs) IDs() []string {
	if jobs == nil {
		return nil
	}
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
	ID        string      `json:"-"`
	JobID     string      `json:"-" db:"job_id"`
	Data      interface{} `json:"data,omitempty"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
}

type JobResult struct {
	ID          string      `json:"-"`
	JobID       string      `json:"-" db:"job_id"`
	Attempt     int         `json:"attempt"`
	Status      *Status     `json:"status"`
	Data        interface{} `json:"data,omitempty"`
	Error       string      `json:"error,omitempty"`
	StartedAt   time.Time   `json:"started_at" db:"started_at"`
	CompletedAt time.Time   `json:"completed_at" db:"completed_at"`
}

type JobListParams struct {
	Names         []string         `json:"names,omitempty"`
	Statuses      []*Status        `json:"statuses,omitempty"`
	Selectors     *label.Selectors `json:"selectors,omitempty"`
	Claims        label.Claims     `json:"claims,omitempty"`
	EnqueuedSince time.Time        `json:"enqueued_since,omitempty"`
	EnqueuedUntil time.Time        `json:"enqueued_until,omitempty"`
	// NoUnclaimed will exclude jobs that have outstanding claims.
	NoUnclaimed bool `json:"no_unclaimed,omitempty"`
}