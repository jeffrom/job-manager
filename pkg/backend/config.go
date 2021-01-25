package backend

import (
	"time"

	"github.com/jeffrom/job-manager/pkg/logger"
)

type Config struct {
	// QueueHistoryLimit is the maximum number of resource versions to store.
	QueueHistoryLimit int           `json:"queue_history_limit" envconfig:"queue_history_limit"`
	ReapAge           time.Duration `json:"reap_age" envconfig:"reap_age"`
	ReapMax           int           `json:"reap_max" envconfig:"reap_max"`
	Debug             bool
	TestMode          bool

	Logger *logger.Logger
}

var DefaultConfig = Config{
	QueueHistoryLimit: 10,
	ReapAge:           10 * time.Minute,
	Debug:             true,
	// ReapAge:           time.Hour * 24 * 60,
}
