package testenv

import (
	"context"
	"testing"

	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/resource"
)

type BackendTestConfig struct {
	Backend backend.Interface

	// Fail fails instead of skipping tests when dependency check (ping) fails.
	Fail bool
}

type backendTestContext struct {
	cfg BackendTestConfig
}

func (tc *backendTestContext) wrap(ctx context.Context, fn func(ctx context.Context, t *testing.T, tc *backendTestContext)) func(t *testing.T) {
	return func(t *testing.T) {
		fn(ctx, t, tc)
	}
}

func BackendTest(cfg BackendTestConfig) func(t *testing.T) {
	tc := &backendTestContext{cfg: cfg}
	return func(t *testing.T) {
		ctx := context.Background()
		be := cfg.Backend

		t.Logf("health checking %T", be)
		if err := be.Ping(ctx); err != nil {
			if cfg.Fail {
				t.Fatalf("backend %T is not healthy", be)
			} else {
				t.Skipf("Skipping backend test for %T because it's not responding to health checks", be)
			}
		}

		t.Logf("resetting %T", be)
		if err := be.Reset(ctx); err != nil {
			t.Fatal(err)
		}

		if !t.Run("queue-admin", tc.wrap(ctx, testQueueAdmin)) {
			return
		}
		if !t.Run("enqueue", tc.wrap(ctx, testEnqueue)) {
			return
		}
	}
}

func testQueueAdmin(ctx context.Context, t *testing.T, tc *backendTestContext) {
	t.Run("save", func(t *testing.T) {
		q, err := tc.cfg.Backend.SaveQueue(ctx, getBasicQueue())
		if err != nil {
			t.Fatal(err)
		}
		if q == nil {
			t.Fatal("queue was nil")
		}
	})
}

func testEnqueue(ctx context.Context, t *testing.T, tc *backendTestContext) {

}

func getBasicQueue() *resource.Queue {
	return &resource.Queue{
		ID:          "cool",
		Concurrency: 3,
		Retries:     3,
	}
}

func getBasicJobs() *resource.Jobs {
	return &resource.Jobs{
		Jobs: []*resource.Job{
			{
				Name: "cooljob",
				Args: jobArgs("nice"),
			},
		},
	}
}

func jobArgs(args ...interface{}) []interface{} { return args }
