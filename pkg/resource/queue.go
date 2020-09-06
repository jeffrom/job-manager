package resource

import (
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
	Unique          bool          `json:"unique,omitempty"`
	Labels          label.Labels  `json:"labels,omitempty"`
	SchemaRaw       []byte        `json:"schema_raw,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at,omitempty"`
	DeletedAt       time.Time     `json:"deleted_at,omitempty"`
}

type Queues struct {
	Queues []*Queue `json:"queues"`
}
