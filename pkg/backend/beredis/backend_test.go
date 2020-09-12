package beredis

import (
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/testenv"
)

func TestBackendRedis(t *testing.T) {
	be := New(WithConfig(Config{
		Config: backend.DefaultConfig,
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
