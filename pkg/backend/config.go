package backend

import "github.com/jeffrom/job-manager/pkg/logger"

type Config struct {
	// HistoryLimit is the maximum number of resource versions to store.
	HistoryLimit int

	Debug    bool
	TestMode bool
	Logger   *logger.Logger
}

var DefaultConfig = Config{
	HistoryLimit: 10,
	Debug:        true,
}
