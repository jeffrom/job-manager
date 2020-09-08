package resource

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jeffrom/job-manager/pkg/label"
)

type Queue struct {
	ID              string        `json:"id"`
	Version         *Version      `json:"version"`
	Concurrency     int           `json:"concurrency,omitempty"`
	Retries         int           `json:"retries,omitempty"`
	Duration        time.Duration `json:"duration,omitempty"`
	CheckinDuration time.Duration `json:"checkin_duration,omitempty"`
	ClaimDuration   time.Duration `json:"claim_duration,omitempty"`
	Unique          bool          `json:"unique,omitempty"`
	Labels          label.Labels  `json:"labels,omitempty"`
	SchemaRaw       []byte        `json:"schema_raw,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at,omitempty"`
	DeletedAt       time.Time     `json:"deleted_at,omitempty"`
}

func (q *Queue) ClaimExpired(job *Job, now time.Time) bool {
	expireTime := job.EnqueuedAt.Add(q.ClaimDuration)
	expired := now.Equal(expireTime) || now.After(expireTime)
	fmt.Println("claim duration:", q.ClaimDuration, "enqueued at:", job.EnqueuedAt.Format(time.Stamp), "now:", now.Format(time.Stamp), "expire:", expireTime.Format(time.Stamp), "expired:", expired)
	return expired
}

func (q *Queue) Equals(other *Queue) bool {
	return q.ID == other.ID &&
		q.Version == other.Version &&
		q.Concurrency == other.Concurrency &&
		q.Retries == other.Retries &&
		q.Duration == other.Duration &&
		q.CheckinDuration == other.CheckinDuration &&
		q.ClaimDuration == other.ClaimDuration &&
		q.Unique == other.Unique &&
		q.Labels.Equals(other.Labels) &&
		bytes.Equal(q.SchemaRaw, other.SchemaRaw)
}

type Queues struct {
	Queues []*Queue `json:"queues"`
}

type QueueListParams struct {
	Names     []string         `json:"names,omitempty"`
	Selectors *label.Selectors `json:"selectors,omitempty"`
}
