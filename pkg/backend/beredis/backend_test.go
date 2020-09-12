package beredis

import (
	"os"
	"testing"

	"github.com/jeffrom/job-manager/pkg/testenv"
)

func TestBackendRedis(t *testing.T) {
	be := New()
	testenv.BackendTest(testenv.BackendTestConfig{
		Backend: be,
		Fail:    os.Getenv("CI") != "",
	})(t)
}
