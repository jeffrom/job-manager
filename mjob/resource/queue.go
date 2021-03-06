package resource

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/jeffrom/job-manager/mjob/label"
)

type Queue struct {
	ID              string       `json:"-"`
	Name            string       `json:"name"`
	Version         *Version     `json:"version" db:"v"`
	Retries         int          `json:"retries,omitempty"`
	Duration        Duration     `json:"duration,omitempty"`
	CheckinDuration Duration     `json:"checkin_duration,omitempty" db:"checkin_duration"`
	ClaimDuration   Duration     `json:"claim_duration,omitempty" db:"claim_duration"`
	Unique          bool         `json:"unique,omitempty" db:"unique_args"`
	Labels          label.Labels `json:"labels,omitempty"`
	SchemaRaw       []byte       `json:"schema_raw,omitempty" db:"job_schema"`
	BackoffInitial  Duration     `json:"backoff_initial" db:"backoff_initial_duration"`
	BackoffMax      Duration     `json:"backoff_max" db:"backoff_max_duration"`
	BackoffFactor   float32      `json:"backoff_factor" db:"backoff_factor"`
	Paused          bool         `json:"paused"`
	Unpaused        bool         `json:"unpaused"`
	Blocked         bool         `json:"blocked"`
	CreatedAt       time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time    `json:"-" db:"updated_at"`
	DeletedAt       NullTime     `json:"deleted_at,omitempty" db:"deleted_at"`
}

func (q *Queue) String() string {
	b, _ := json.Marshal(q)
	return string(b)
}

func (q *Queue) ClaimExpired(job *Job, now time.Time) bool {
	lastClaim := job.LastClaimWindow()
	expireTime := lastClaim.Add(time.Duration(q.ClaimDuration))
	expired := now.Equal(expireTime) || now.After(expireTime)
	// fmt.Println("claim duration:", q.ClaimDuration, "last claim:", lastClaim.Format(time.Stamp), "now:", now.Format(time.Stamp), "expire:", expireTime.Format(time.Stamp), "expired:", expired)
	return expired
}

func (q *Queue) EqualAttrs(other *Queue) bool {
	// fmt.Printf("a: %+v\nb: %+v\n", q, other)
	return q.Name == other.Name &&
		q.Retries == other.Retries &&
		q.Duration == other.Duration &&
		q.CheckinDuration == other.CheckinDuration &&
		q.ClaimDuration == other.ClaimDuration &&
		q.Unique == other.Unique &&
		fmt.Sprint(q.BackoffFactor) == fmt.Sprint(other.BackoffFactor) &&
		q.BackoffInitial == other.BackoffInitial &&
		q.BackoffMax == other.BackoffMax &&
		q.Labels.Equals(other.Labels) &&
		// (q.DeletedAt == nil || q.DeletedAt.Valid == false) == (other.DeletedAt == nil || other.DeletedAt.Valid == false) &&
		q.DeletedAt.Valid == other.DeletedAt.Valid &&
		jsonEquals(q.SchemaRaw, other.SchemaRaw)
}

func (q *Queue) Equal(other *Queue) bool {
	return q.CreatedAt.Equal(other.CreatedAt) &&
		q.UpdatedAt.Equal(other.UpdatedAt) &&
		q.DeletedAt.Valid == other.DeletedAt.Valid &&
		q.EqualAttrs(other)
}

func (q *Queue) Copy() *Queue {
	cp := &Queue{}
	*cp = *q
	if q.Version != nil {
		*cp.Version = *q.Version
	}

	// TODO copy schema bytes?
	// if len(q.SchemaRaw) > 0 {
	// 	copy(cp.SchemaRaw, q.SchemaRaw)
	// }
	return cp
}

type Queues struct {
	Queues []*Queue `json:"queues"`
}

func (qs *Queues) ToMap() map[string]*Queue {
	m := make(map[string]*Queue)
	for _, q := range qs.Queues {
		m[q.Name] = q
	}
	return m
}

type QueueListParams struct {
	Names     []string         `json:"names,omitempty"`
	Selectors *label.Selectors `json:"selectors,omitempty"`
	Page      *Pagination      `json:"page,omitempty"`
	Includes  []string         `json:"include,omitempty"`
}

func jsonEquals(a, b []byte) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	var ai interface{}
	var bi interface{}
	aerr := json.Unmarshal(a, &ai)
	berr := json.Unmarshal(b, &bi)
	if aerr != berr {
		return false
	}

	return reflect.DeepEqual(ai, bi)
}
