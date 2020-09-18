package bepg

import (
	"strconv"
	"strings"

	"github.com/jeffrom/job-manager/pkg/backend"
)

type Config struct {
	backend.Config

	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	SSLMode  string `json:"sslmode,omitempty"`
}

func (c Config) DSN() string {
	var b strings.Builder

	n := 0
	if c.Host != "" {
		semanticSpace(b, n)
		b.WriteString("host=")
		b.WriteString(c.Host)
		n++
	}

	if c.Port > 0 {
		semanticSpace(b, n)
		b.WriteString("port=")
		b.WriteString(strconv.FormatInt(int64(c.Port), 10))
		n++
	}

	if c.User != "" {
		semanticSpace(b, n)
		b.WriteString("user=")
		b.WriteString(c.User)
		n++
	}

	if c.Password != "" {
		semanticSpace(b, n)
		b.WriteString("password=")
		b.WriteString(c.Password)
		n++
	}

	if c.SSLMode != "" {
		semanticSpace(b, n)
		b.WriteString("sslmode=")
		b.WriteString(c.SSLMode)
		n++
	}

	return b.String()
}

func semanticSpace(b strings.Builder, n int) {
	if n > 0 {
		b.WriteByte(' ')
	}
}
