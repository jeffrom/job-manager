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

		ctx = mustReset(ctx, t, be)
		if !t.Run("queue-admin", tc.wrap(ctx, testQueueAdmin)) {
			return
		}
		if !t.Run("enqueue-dequeue", tc.wrap(ctx, testEnqueueDequeue)) {
			return
		}

		// ctx = mustReset(ctx, t, be)
	}
}

var basictime = time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)

func testQueueAdmin(ctx context.Context, t *testing.T, tc *backendTestContext) {
	t.Run("save", func(t *testing.T) {
		be := tc.cfg.Backend

		// initial creation should work
		now := basictime
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

	t.Run("get", func(t *testing.T) {

	})

	t.Run("delete", func(t *testing.T) {

	})
}

func testEnqueueDequeue(ctx context.Context, t *testing.T, tc *backendTestContext) {
	be := tc.cfg.Backend
	ctx = mustReset(ctx, t, be)

	now := basictime
	ctx = internal.SetMockTime(ctx, now)
	mustSaveQueue(ctx, t, be, getBasicQueue())

	expectJobs := getBasicJobs()
	res := mustEnqueueJobs(ctx, t, be, expectJobs)

	jobs := res.Jobs
	if len(jobs) != 3 {
		t.Fatalf("expected 3 jobs, got %d", len(jobs))
	}
	checkJob(t, jobs[0])
	checkJobStatus(t, resource.StatusQueued, jobs[0])
	checkVersion(t, 1, jobs[0].Version)
	checkJob(t, jobs[1])
	checkJobStatus(t, resource.StatusQueued, jobs[1])
	checkVersion(t, 1, jobs[1].Version)
	checkJob(t, jobs[2])
	checkJobStatus(t, resource.StatusQueued, jobs[2])
	checkVersion(t, 1, jobs[2].Version)

	now = now.Add(1 * time.Second)
	ctx = internal.SetMockTime(ctx, now)
	// now dequeue them
	deqRes, err := be.DequeueJobs(ctx, 3, &resource.JobListParams{
		Names: []string{"cool"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if deqRes == nil {
		t.Fatal("expected dequeue result not to be nil")
	}
	if l := len(deqRes.Jobs); l != 3 {
		t.Fatalf("expected to dequeue 3 jobs, got %d", l)
	}
	deqJobs := deqRes.Jobs
	checkJob(t, deqJobs[0])
	checkJobStatus(t, resource.StatusRunning, deqJobs[0])
	checkVersion(t, 2, deqJobs[0].Version)
	checkJob(t, deqJobs[1])
	checkJobStatus(t, resource.StatusRunning, deqJobs[1])
	checkVersion(t, 2, deqJobs[1].Version)
	checkJob(t, deqJobs[2])
	checkJobStatus(t, resource.StatusRunning, deqJobs[2])
	checkVersion(t, 2, deqJobs[2].Version)
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
				Name: "cool",
				Args: jobArgs(1, "nice"),
			},
			{
				Name: "cool",
				Args: jobArgs(2, "222"),
			},
			{
				Name: "cool",
				Args: jobArgs(3, "heck"),
			},
		},
	}
}

func jobArgs(args ...interface{}) []interface{} { return args }

func mustReset(ctx context.Context, t testing.TB, be backend.Interface) context.Context {
	t.Helper()
	t.Logf("Resetting %T", be)
	if err := be.Reset(ctx); err != nil {
		t.Fatal(err)
	}

	return internal.SetTimeProvider(ctx, internal.Time(0))
}

func mustSaveQueue(ctx context.Context, t testing.TB, be backend.Interface, q *resource.Queue) *resource.Queue {
	t.Helper()

	t.Logf("SaveQueue(%+v)", readable(q))
	res, err := be.SaveQueue(ctx, q)
	if err != nil {
		t.Logf("-> err: %v", err)
		t.Fatal(err)
	}
	t.Logf("-> %s", readable(res))
	return res
}

func mustEnqueueJobs(ctx context.Context, t testing.TB, be backend.Interface, jobs *resource.Jobs) *resource.Jobs {
	t.Helper()

	t.Logf("EnqueueJobs(%+v)", readable(jobs.Jobs))
	res, err := be.EnqueueJobs(ctx, jobs)
	if err != nil {
		t.Fatal(err)
	}

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

func checkJob(t testing.TB, jb *resource.Job) bool {
	t.Helper()

	if jb == nil {
		t.Fatal("expected job not to be nil")
	}
	if jb.ID == "" {
		t.Error("expected job id to be set")
	}
	if v := jb.Version; v == nil || v.Raw() < 1 {
		t.Error("expected job version to be > v1, was", v)
	}
	if v := jb.QueueVersion; v == nil || v.Raw() < 1 {
		t.Error("expected job queue version to be > v1, was", v)
	}

	if jb.EnqueuedAt.IsZero() {
		t.Error("expected job enqueued_at to be set")
	}
	if loc := jb.EnqueuedAt.Location(); loc == nil || loc != time.UTC {
		t.Errorf("queue created_at timezone was not utc, was %q", loc.String())
	}

	return !t.Failed()
}

func checkJobStatus(t testing.TB, expect resource.Status, jb *resource.Job) bool {
	t.Helper()
	if st := jb.Status; st != expect {
		t.Errorf("expected job status %s, got %s", expect, st)
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