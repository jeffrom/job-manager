package mjob

import "github.com/jeffrom/job-manager/mjob/client"

type ConsumerConfig struct {
	DequeueOpts  client.DequeueOpts `json:"dequeue_opts"`
	Concurrency  int                `json:"concurrency"`
	Backpressure int                `json:"backpressure"`
}

var defaultConsumerConfig ConsumerConfig = ConsumerConfig{
	Concurrency: 1,
}
