package integration

import (
	"testing"

	"github.com/jeffrom/job-manager/pkg/resource"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
)

func checkArgsSchema(t testing.TB, verr *resource.ValidationError, expectPath string) {
	t.Helper()
	if verr.Path != expectPath {
		t.Errorf("expected path %q, got %q", expectPath, verr.Path)
	}
}

func checkSchemaNumber(t testing.TB, v interface{}, expect float64) {
	t.Helper()
	val, ok := v.(float64)
	if !ok {
		t.Fatalf("expected invalid value type float64, got %T", v)
	}
	if val != expect {
		t.Errorf("expected value %f, got %f", expect, val)
	}
}

func checkSchemaString(t testing.TB, v interface{}, expect string) {
	t.Helper()
	val, ok := v.(string)
	if !ok {
		t.Fatalf("expected invalid value type string, got %T", v)
	}
	if val != expect {
		t.Errorf("expected value %q, got %q", expect, val)
	}
}

func checkJob(t testing.TB, jobData *jobv1.Job) {
	if jobData.EnqueuedAt == nil {
		t.Errorf("jobv1.EnqueuedAt was nil")
	} else if !jobData.EnqueuedAt.IsValid() {
		t.Errorf("jobv1.EnqueuedAt is invalid: %v", jobData.EnqueuedAt.CheckValid())
	} else if jobData.EnqueuedAt.AsTime().IsZero() {
		t.Errorf("jobv1.EnqueuedAt is zero")
	}

	if jobData.Status == jobv1.StatusUnknown {
		t.Errorf("jobv1.Status is %s", jobv1.StatusUnknown)
	}
	if jobData.Attempt < 0 {
		t.Errorf("jobv1.Attempt was %d", jobData.Attempt)
	}
}
