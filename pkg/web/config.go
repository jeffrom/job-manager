package web

import (
	"io"
	"os"

	"github.com/jeffrom/job-manager/pkg/backend"
)

type Config struct {
	LogOutput io.Writer
	Backend   string `json:"backend"`
	be        backend.Interface
}

func NewConfig() Config {
	return Config{
		LogOutput: os.Stdout,
		Backend:   "memory",
	}
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
