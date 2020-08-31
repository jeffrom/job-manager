package middleware

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/jeffrom/job-manager/pkg/backend"
)

func init() {
	zerolog.TimestampFieldName = "timestamp"
	zerolog.MessageFieldName = "msg"
	zerolog.DurationFieldUnit = time.Millisecond
	// zerolog.ErrorMarshalFunc = logErrorMarshaler
}

type Config struct {
	LogJSON         bool          `json:"log_json"`
	Backend         string        `json:"backend"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`

	DefaultMaxJobTimeout time.Duration `json:"default_max_job_timeout"`
	DefaultConcurrency   int           `json:"json:"default_concurrency""`
	DefaultMaxRetries    int           `json:"default_max_retries"`

	Logger    zerolog.Logger `json:"-"`
	LogOutput io.Writer      `json:"-"`
	be        backend.Interface
}

func NewConfig() Config {
	out := os.Stdout
	c := Config{
		LogOutput:            out,
		Backend:              "memory",
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
}

func (c *Config) GetBackend() backend.Interface {
	if c.be != nil {
		return c.be
	}
	switch c.Backend {
	case "":
		return nil
	case "memory":
		c.be = backend.NewMemory()
		return c.be
	default:
		panic("unsupported backend: " + c.Backend)
	}
}

func (c *Config) newLogger(out io.Writer) zerolog.Logger {
	if !c.LogJSON {
		out = zerolog.ConsoleWriter{Out: out}
	}
	c.LogOutput = out
	l := zerolog.New(out).With().Timestamp().Logger()

	return l
}

func ConfigFromContext(ctx context.Context) Config {
	return ctx.Value(ConfigKey).(Config)
}
