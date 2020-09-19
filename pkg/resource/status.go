package resource

import (
	"database/sql/driver"
	"strconv"
)

type Status int

const (
	StatusUnspecified Status = iota
	StatusQueued
	StatusRunning
	StatusComplete
	StatusDead
	StatusCancelled
	StatusInvalid
	StatusFailed
)

func NewStatus(s Status) *Status { return &s }

func statusFromString(s string) Status {
	switch s {
	case "unspecified":
		return StatusUnspecified
	case "queued":
		return StatusQueued
	case "running":
		return StatusRunning
	case "complete":
		return StatusComplete
	case "failed":
		return StatusFailed
	case "dead":
		return StatusDead
	case "invalid":
		return StatusInvalid
	case "cancelled":
		return StatusCancelled
	default:
		return StatusUnspecified
	}
}

func (s *Status) String() string {
	switch *s {
	case StatusUnspecified:
		return "unspecified"
	case StatusQueued:
		return "queued"
	case StatusRunning:
		return "running"
	case StatusComplete:
		return "complete"
	case StatusFailed:
		return "failed"
	case StatusDead:
		return "dead"
	case StatusInvalid:
		return "invalid"
	case StatusCancelled:
		return "cancelled"
	default:
		panic("unknown status: " + strconv.FormatInt(int64(*s), 10))
	}
}

func (s *Status) Scan(value interface{}) error {
	if value == nil {
		*s = StatusUnspecified
		return nil
	}

	valstr := value.(string)
	*s = statusFromString(valstr)
	return nil
}

func (s *Status) Value() (driver.Value, error) {
	return s.String(), nil
}

func StatusIsAttempted(status *Status) bool {
	switch *status {
	case StatusComplete, StatusCancelled, StatusInvalid, StatusDead:
		return true
	}
	return false
}
