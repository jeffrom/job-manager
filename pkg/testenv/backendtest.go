package testenv

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/internal"
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
		be := tc.cfg.Backend

		// initial creation should work
		now := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
		ctx = internal.SetMockTime(ctx, now)
		q := mustSaveQueue(ctx, t, be, getBasicQueue())
		mustCheck(t, checkQueue(t, q))
		mustCheck(t, checkVersion(t, 1, q.Version))

		// same attrs should result in no changes
		nextNow := now.Add(1 * time.Second)
		ctx = internal.SetMockTime(ctx, nextNow)
		q = mustSaveQueue(ctx, t, be, getBasicQueue())
		mustCheck(t, checkQueue(t, q))
		mustCheck(t, checkVersion(t, 1, q.Version))
		if !q.CreatedAt.Equal(now) {
			t.Fatalf("expected created_at to be %q, was %q", now, q.CreatedAt)
		}
		if !q.UpdatedAt.Equal(now) {
			t.Fatalf("expected updated_at to be %q, was %q", now, q.UpdatedAt)
		}

		// a legit update should work
		lastNow := now.Add(1 * time.Second)
		ctx = internal.SetMockTime(ctx, lastNow)
		q2 := getBasicQueue()
		q2.Concurrency++
		res := mustSaveQueue(ctx, t, be, q2)
		mustCheck(t, checkQueue(t, res))
		mustCheck(t, checkVersion(t, 2, res.Version))
		if !res.CreatedAt.Equal(now) {
			t.Fatalf("expected created_at to be %q, was %q", now, res.CreatedAt)
		}
		if !res.UpdatedAt.Equal(lastNow) {
			t.Fatalf("expected updated_at to be %q, was %q", lastNow, res.UpdatedAt)
		}
	})
}

func testEnqueue(ctx context.Context, t *testing.T, tc *backendTestContext) {

}

func getBasicQueue() *resource.Queue {
	return &resource.Queue{
		ID:          "cool",
		Version:     resource.NewVersion(1),
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

func mustSaveQueue(ctx context.Context, t testing.TB, be backend.Interface, q *resource.Queue) *resource.Queue {
	t.Helper()

	res, err := be.SaveQueue(ctx, q)
	t.Logf("SaveQueue(%+v)", readable(q))
	if err != nil {
		t.Logf("-> err: %v", err)
		t.Fatal(err)
	}
	t.Logf("-> %s", readable(res))
	return res
}

func mustCheck(t testing.TB, res bool) {
	t.Helper()
	if t.Failed() {
		t.FailNow()
	}
}

func checkQueue(t testing.TB, q *resource.Queue) bool {
	t.Helper()
	if q == nil {
		t.Error("queue was nil")
		return !t.Failed()
	}

	if q.ID == "" {
		t.Error("queue id was empty")
	}

	if loc := q.CreatedAt.Location(); loc == nil || loc != time.UTC {
		t.Errorf("queue created_at timezone was not utc, was %q", loc.String())
	}
	if loc := q.UpdatedAt.Location(); loc == nil || loc != time.UTC {
		t.Errorf("queue updated_at timezone was not utc, was %q", loc.String())
	}
	if q.Version.Raw() == 1 {
		if !q.CreatedAt.Equal(q.UpdatedAt) {
			t.Errorf("expected new queue updated_at to be set by backend to equal created_at (created_at: %q, updated_at: %q)", q.CreatedAt, q.UpdatedAt)
		}
	} else {
		if !q.UpdatedAt.After(q.CreatedAt) {
			t.Errorf("expected updated queue updated_at to be later than created_at (created_at: %q, updated_at: %q)", q.CreatedAt, q.UpdatedAt)
		}
	}

	if v := q.Version; v == nil {
		t.Error("queue version was nil")
	} else if v.Raw() == 1 {

	}
	return !t.Failed()
}

func checkVersion(t testing.TB, expect int32, v *resource.Version) bool {
	t.Helper()

	if v == nil {
		t.Error("version was nil")
		return !t.Failed()
	}

	if raw := v.Raw(); raw != expect {
		t.Errorf("expected version to be v%d, was %s", expect, v)
	}

	return !t.Failed()
}

var jsonRE = regexp.MustCompile(`"(\w+)":`)

func readable(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {

	}
	s := string(b)
	if len(s) == 0 {
		return s
	}

	if s[0] == '{' && s[len(s)-1] == '}' {
		s = s[1 : len(s)-1]
	}
	return strings.ReplaceAll(
		strings.TrimSpace(jsonRE.ReplaceAllString(s, " $1: ")),
		`"0001-01-01T00:00:00Z"`, "0")
}
