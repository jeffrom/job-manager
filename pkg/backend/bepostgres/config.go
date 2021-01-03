package bepostgres

import (
	"strconv"
	"strings"

	"github.com/jeffrom/job-manager/pkg/backend"
)

var DefaultConfig = Config{
	Config:   backend.DefaultConfig,
	Database: "job_manager",
	// Host:     "localhost",
}

type Config struct {
	backend.Config

	Database string `json:"database" envconfig:"postgres_database"`
	Host     string `json:"host,omitempty" envconfig:"postgres_host"`
	Port     int    `json:"port,omitempty" envconfig:"postgres_port"`
	User     string `json:"user,omitempty" envconfig:"postgres_user"`
	Password string `json:"password,omitempty" envconfig:"postgres_pass"`
	SSLMode  string `json:"sslmode,omitempty" envconfig:"postgres_sslmode"`
}

func (c Config) DSN() string {
	var b strings.Builder

	n := 0
	if c.Database != "" {
		semanticSpace(&b, n)
		b.WriteString("database=")
		b.WriteString(c.Database)
		n++
	}

	if c.Host != "" {
		semanticSpace(&b, n)
		b.WriteString("host=")
		b.WriteString(c.Host)
		n++
	}

	if c.Port > 0 {
		semanticSpace(&b, n)
		b.WriteString("port=")
		b.WriteString(strconv.FormatInt(int64(c.Port), 10))
		n++
	}

	if c.User != "" {
		semanticSpace(&b, n)
		b.WriteString("user=")
		b.WriteString(c.User)
		n++
	}

	if c.Password != "" {
		semanticSpace(&b, n)
		b.WriteString("password=")
		b.WriteString(c.Password)
		n++
	}

	if c.SSLMode != "" {
		semanticSpace(&b, n)
		b.WriteString("sslmode=")
		b.WriteString(c.SSLMode)
		n++
	}

	return b.String()
}

func semanticSpace(b *strings.Builder, n int) {
	if n > 0 {
		b.WriteByte(' ')
	}
}
