// Package mjob is the client entrypoint package for job-manager.
package mjob

import (
	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/consumer"
	"github.com/jeffrom/job-manager/mjob/resource"
)

// Job is the job type declared in mjob/resource.
type Job = resource.Job

// Result is the job result type declared in mjob/resource.
type Result = resource.JobResult

// Client is the interface for the job-manager http client in package
// mjob/client.
type Client = client.Interface

// NewConsumer creates a new job consumer. Consumers can target one or all
// queues. See mjob/consumer.Provider docs for confuration options.
func NewConsumer(c client.Interface, r consumer.Runner, providers ...consumer.Provider) *consumer.Consumer {
	return consumer.New(c, r, providers...)
}
