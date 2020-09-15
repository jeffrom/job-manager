package resource

import (
	"bytes"
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
	lastClaim := job.LastClaimWindow()
	expireTime := lastClaim.Add(q.ClaimDuration)
	expired := now.Equal(expireTime) || now.After(expireTime)
	// fmt.Println("claim duration:", q.ClaimDuration, "last claim:", lastClaim.Format(time.Stamp), "now:", now.Format(time.Stamp), "expire:", expireTime.Format(time.Stamp), "expired:", expired)
	return expired
}

func (q *Queue) EqualAttrs(other *Queue) bool {
	return q.ID == other.ID &&
		q.Concurrency == other.Concurrency &&
		q.Retries == other.Retries &&
		q.Duration == other.Duration &&
		q.CheckinDuration == other.CheckinDuration &&
		q.ClaimDuration == other.ClaimDuration &&
		q.Unique == other.Unique &&
		q.Labels.Equals(other.Labels) &&
		q.DeletedAt.IsZero() == other.DeletedAt.IsZero() &&
		bytes.Equal(q.SchemaRaw, other.SchemaRaw)
}

func (q *Queue) Equal(other *Queue) bool {
	return q.CreatedAt.Equal(other.CreatedAt) &&
		q.UpdatedAt.Equal(other.UpdatedAt) &&
		q.DeletedAt.Equal(other.DeletedAt) &&
		q.EqualAttrs(other)
}

func (q *Queue) Copy() *Queue {
	cp := &Queue{}
	*cp = *q
	if q.Version != nil {
		*cp.Version = *q.Version
	}
	return cp
}

type Queues struct {
	Queues []*Queue `json:"queues"`
}

func (qs *Queues) ToMap() map[string]*Queue {
	m := make(map[string]*Queue)
	for _, q := range qs.Queues {
		m[q.ID] = q
	}
	return m
}

type QueueListParams struct {
	Names     []string         `json:"names,omitempty"`
	Selectors *label.Selectors `json:"selectors,omitempty"`
}
