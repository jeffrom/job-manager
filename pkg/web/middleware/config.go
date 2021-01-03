package middleware

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/jeffrom/job-manager/pkg/logger"
)

func init() {
	zerolog.TimestampFieldName = "timestamp"
	zerolog.MessageFieldName = "msg"
	zerolog.DurationFieldUnit = time.Millisecond
	// zerolog.ErrorMarshalFunc = logErrorMarshaler
}

type Config struct {
	Host            string        `json:"host"`
	LogJSON         bool          `json:"log_json"`
	DebugLog        bool          `json:"debug_log"`
	Backend         string        `json:"backend"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`

	DefaultMaxJobTimeout time.Duration `json:"default_max_job_timeout"`
	DefaultConcurrency   int           `json:"json:"default_concurrency""`
	DefaultMaxRetries    int           `json:"default_max_retries"`

	Logger    *logger.Logger `json:"-"`
	LogOutput io.Writer      `json:"-"`
}

var ConfigDefaults = Config{
	Host:                 ":1874",
	DebugLog:             true,
	Backend:              "postgres",
	ShutdownTimeout:      30 * time.Second,
	DefaultMaxJobTimeout: 10 * time.Minute,
	DefaultConcurrency:   10,
	DefaultMaxRetries:    10,
}

func NewConfig() Config {
	out := os.Stdout
	c := Config{
		Host:                 ":1874",
		LogOutput:            out,
		DebugLog:             true,
		Backend:              "postgres",
		ShutdownTimeout:      30 * time.Second,
		DefaultMaxJobTimeout: 10 * time.Minute,
		DefaultConcurrency:   10,
		DefaultMaxRetries:    10,
	}
	c.Logger = c.newLogger(out)
	return c
}

func (c *Config) ResetLogOutput(out io.Writer) {
	c.Logger = c.newLogger(out)
	c.LogOutput = out
}

func (c *Config) newLogger(out io.Writer) *logger.Logger {
	return logger.New(out, c.LogJSON)
}

func ConfigFromContext(ctx context.Context) Config {
	return ctx.Value(ConfigKey).(Config)
}
