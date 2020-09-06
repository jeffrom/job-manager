package resource

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

func StatusIsAttempted(status Status) bool {
	switch status {
	case StatusComplete, StatusCancelled, StatusInvalid, StatusDead:
		return true
	}
	return false
}
