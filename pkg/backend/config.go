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
	TestMode          bool          `json:"test_mode" envconfig:"test"`
	Debug             bool

	Logger *logger.Logger
}

var DefaultConfig = Config{
	QueueHistoryLimit: 10,
	ReapAge:           24 * time.Hour * 60,
	// Debug:             true,
	// ReapAge:           10 * time.Minute,
}
