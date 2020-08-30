// Package job contains the protocol for working with jobs and queues.
package job

import uuid "github.com/satori/go.uuid"

func NewID() string {
	return uuid.NewV4().String()
}

func HasStatus(job *Job, statuses []Status) bool {
	for _, st := range statuses {
		if job.Status == st {
			return true
		}
	}
	return false
}
