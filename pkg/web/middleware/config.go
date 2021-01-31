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
	DevLog          bool          `json:"dev_log" envconfig:"dev_log"`
	DebugLog        bool          `json:"debug_log" envconfig:"debug"`
	Backend         string        `json:"backend" envconfig:"backend"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout" envconfig:"shutdown_timeout"`

	DefaultMaxJobTimeout time.Duration `json:"default_max_job_timeout" envconfig:"default_max_job_timeout"`
	DefaultMaxRetries    int           `json:"default_max_retries" envconfig:"default_max_retries"`

	InvalidateInterval time.Duration `json:"invalidate_interval,omitempty" envconfig:"invalidate_interval"`
	ReapInterval       time.Duration `json:"reap_interval,omitempty" envconfig:"reap_interval"`
	ReapAge            time.Duration `json:"reap_age,omitempty" envconfig:"reap_age"`
	ReapMax            int           `json:"reap_max,omitempty" envconfig:"reap_max"`

	Logger    *logger.Logger `json:"-"`
	LogOutput io.Writer      `json:"-"`
}

var ConfigDefaults = Config{
	Host:                 ":1874",
	Backend:              "postgres",
	ShutdownTimeout:      30 * time.Second,
	DefaultMaxJobTimeout: 10 * time.Minute,
	DefaultMaxRetries:    10,
	InvalidateInterval:   15 * time.Second,
	ReapInterval:         10 * time.Minute,
	// DebugLog:             true,
}

func NewConfig() Config {
	out := os.Stdout
	c := Config{
		Host:                 ":1874",
		LogOutput:            out,
		Backend:              "postgres",
		ShutdownTimeout:      30 * time.Second,
		DefaultMaxJobTimeout: 10 * time.Minute,
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
	return logger.New(out, !c.DevLog, c.DebugLog)
}

func ConfigFromContext(ctx context.Context) Config {
	return ctx.Value(ConfigKey).(Config)
}
