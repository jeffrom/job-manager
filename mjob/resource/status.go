package resource

import (
	"database/sql/driver"
	"encoding/json"
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

func StatusFromString(s string) *Status {
	switch s {
	case "unspecified":
		return NewStatus(StatusUnspecified)
	case "queued":
		return NewStatus(StatusQueued)
	case "running":
		return NewStatus(StatusRunning)
	case "complete":
		return NewStatus(StatusComplete)
	case "failed":
		return NewStatus(StatusFailed)
	case "dead":
		return NewStatus(StatusDead)
	case "invalid":
		return NewStatus(StatusInvalid)
	case "cancelled":
		return NewStatus(StatusCancelled)
	default:
		return NewStatus(StatusUnspecified)
	}
}

func StatusesFromStrings(statuses ...string) []*Status {
	res := make([]*Status, len(statuses))
	for i, st := range statuses {
		res[i] = StatusFromString(st)
	}
	return res
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
	*s = *StatusFromString(valstr)
	return nil
}

func (s *Status) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var statusStr string
	if err := json.Unmarshal(data, statusStr); err != nil {
		return err
	}
	st := StatusFromString(statusStr)
	*s = *st
	return nil
}

func StatusIsDone(status *Status) bool {
	switch *status {
	case StatusComplete, StatusCancelled, StatusInvalid, StatusDead:
		return true
	}
	return false
}

func StatusStrings(statuses ...*Status) []string {
	res := make([]string, len(statuses))
	for i, st := range statuses {
		res[i] = st.String()
	}
	return res
}
