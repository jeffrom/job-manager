package consumer

import (
	"time"

	"github.com/jeffrom/job-manager/mjob/client"
)

type Config struct {
	DequeueOpts     client.DequeueOpts `json:"dequeue_opts"`
	Concurrency     int                `json:"concurrency"`
	ShutdownTimeout time.Duration      `json:"shutdown_timeout"`
}

var defaultConfig Config = Config{
	Concurrency:     1,
	ShutdownTimeout: 15 * time.Second,
}
