package web

import "github.com/jeffrom/job-manager/pkg/backend"

type Config struct {
	Backend string `json:"backend"`
	be      backend.Interface
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
