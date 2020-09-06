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
