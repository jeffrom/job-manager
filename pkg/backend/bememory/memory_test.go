package bememory

import (
	"testing"

	"github.com/jeffrom/job-manager/pkg/testenv"
)

func TestBackendMemory(t *testing.T) {
	be := New()
	testenv.BackendTest(testenv.BackendTestConfig{
		Backend: be,
		Fail:    true,
	})(t)
}
