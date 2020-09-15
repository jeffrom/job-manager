package resource

import "strconv"

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

func (s Status) String() string {
	switch s {
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
		panic("unknown status: " + strconv.FormatInt(int64(s), 10))
	}
}

func StatusIsAttempted(status Status) bool {
	switch status {
	case StatusComplete, StatusCancelled, StatusInvalid, StatusDead:
		return true
	}
	return false
}
