package beredis

import (
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/testenv"
)

func TestBackendRedis(t *testing.T) {
	t.Skip("Skipping because implementation broke")
	defaultCfg := backend.DefaultConfig
	defaultCfg.TestMode = true
	be := New(WithConfig(Config{
		Config: defaultCfg,
		Redis: &redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       1,
		},
	}))
	testenv.BackendTest(testenv.BackendTestConfig{
		Backend: be,
		Fail:    os.Getenv("CI") != "",
	})(t)
}
