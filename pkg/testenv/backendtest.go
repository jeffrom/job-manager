package testenv

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/jeffrom/job-manager/mjob/resource"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/internal"
)

type BackendTestConfig struct {
	Type    string
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
		if !t.Run("attempts", tc.wrap(ctx, testAttempts)) {
			return
		}
		if !t.Run("dequeue-while-running", tc.wrap(ctx, testDequeueWhileRunning)) {
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
		q = mustSaveQueue(ctx, t, be, q)
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
		q2 := q
		q2.Retries++
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

	// t.Run("get", func(t *testing.T) {

	// })

	t.Run("delete", func(t *testing.T) {
		be := tc.cfg.Backend
		ctx = mustReset(ctx, t, be)
		q := mustSaveQueue(ctx, t, be, getBasicQueue())
		mustGetQueue(ctx, t, be, q.Name, nil)
		mustDeleteQueue(ctx, t, be, q.Name)
		checkQueueNotFound(ctx, t, be, q.Name)
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

	// now dequeue them
	now = now.Add(1 * time.Second)
	ctx = internal.SetMockTime(ctx, now)
	deqRes := mustDequeueJobs(ctx, t, be, 3, &resource.JobListParams{
		Queues: []string{"cool"},
	})
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

	// now ack
	now = now.Add(1 * time.Second)
	ctx = internal.SetMockTime(ctx, now)
	acks := []*resource.Ack{
		{
			JobID:  jobs[0].ID,
			Status: resource.NewStatus(resource.StatusComplete),
		},
		{
			JobID:  jobs[1].ID,
			Status: resource.NewStatus(resource.StatusComplete),
		},
		{
			JobID:  jobs[2].ID,
			Status: resource.NewStatus(resource.StatusComplete),
		},
	}
	mustAckJobs(ctx, t, be, acks)

	// time.Sleep(1 * time.Second)
	resJobs := mustListJobs(ctx, t, be, 3, &resource.JobListParams{
		Queues:   []string{"cool"},
		Statuses: []*resource.Status{resource.NewStatus(resource.StatusComplete)},
	})

	// this can be eventually consistent, but we should get either 0 or 3 rows back now.
	if l := len(resJobs.Jobs); l == 3 {
		ackedJobs := resJobs.Jobs
		checkJob(t, ackedJobs[0])
		checkJobStatus(t, resource.StatusComplete, ackedJobs[0])
		checkVersion(t, 3, ackedJobs[0].Version)
		checkJob(t, ackedJobs[1])
		checkJobStatus(t, resource.StatusComplete, ackedJobs[1])
		checkVersion(t, 3, ackedJobs[1].Version)
		checkJob(t, ackedJobs[2])
		checkJobStatus(t, resource.StatusComplete, ackedJobs[2])
		checkVersion(t, 3, ackedJobs[2].Version)
	} else if l != 0 {
		t.Fatalf("expected 3 or 0 rows, got %d", l)
	}
}

func testAttempts(ctx context.Context, t *testing.T, tc *backendTestContext) {
	be := tc.cfg.Backend
	if tc.cfg.Type == "memory" {
		t.Skip("skipping as memory backend doesn't currently support attempts")
	}
	ctx = mustReset(ctx, t, be)

	now := basictime
	ctx = internal.SetMockTime(ctx, now)
	q := &resource.Queue{
		Name:           "cool",
		Retries:        2,
		BackoffInitial: resource.Duration(10 * time.Second),
		BackoffFactor:  2.0,
		BackoffMax:     resource.Duration(10 * time.Minute),
	}
	mustSaveQueue(ctx, t, be, q)

	jobs := getBasicJobs()
	jobs.Jobs = []*resource.Job{jobs.Jobs[0]}
	res := mustEnqueueJobs(ctx, t, be, jobs)
	if len(res.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(res.Jobs))
	}

	now = now.Add(1 * time.Second)
	ctx = internal.SetMockTime(ctx, now)
	deqRes := mustDequeueJobs(ctx, t, be, 1, &resource.JobListParams{
		Queues: []string{"cool"},
	})
	if deqRes == nil {
		t.Fatal("expected dequeue result not to be nil")
	}
	if l := len(deqRes.Jobs); l != 1 {
		t.Fatalf("expected to dequeue 1 job, got %d", l)
	}

	now = now.Add(1 * time.Second)
	ctx = internal.SetMockTime(ctx, now)
	acks := []*resource.Ack{
		{
			JobID:  deqRes.Jobs[0].ID,
			Status: resource.NewStatus(resource.StatusFailed),
		},
	}
	mustAckJobs(ctx, t, be, acks)

	now = now.Add(1 * time.Second)
	ctx = internal.SetMockTime(ctx, now)
	deqRes = mustDequeueJobs(ctx, t, be, 1, &resource.JobListParams{
		Queues: []string{"cool"},
	})
	if deqRes == nil {
		t.Fatal("expected dequeue result not to be nil")
	}
	if l := len(deqRes.Jobs); l != 0 {
		t.Fatalf("expected to dequeue 0 jobs, got %d", l)
	}

	now = now.Add(10 * time.Second)
	ctx = internal.SetMockTime(ctx, now)
	deqRes = mustDequeueJobs(ctx, t, be, 1, &resource.JobListParams{
		Queues: []string{"cool"},
	})
	if deqRes == nil {
		t.Fatal("expected dequeue result not to be nil")
	}
	if l := len(deqRes.Jobs); l != 1 {
		t.Fatalf("expected to dequeue 1 job, got %d", l)
	}
	// fail again
	now = now.Add(1 * time.Second)
	ctx = internal.SetMockTime(ctx, now)
	mustAckJobs(ctx, t, be, acks)

	now = now.Add(1 * time.Second)
	ctx = internal.SetMockTime(ctx, now)
	deqRes = mustDequeueJobs(ctx, t, be, 1, &resource.JobListParams{
		Queues: []string{"cool"},
	})
	if deqRes == nil {
		t.Fatal("expected dequeue result not to be nil")
	}
	if l := len(deqRes.Jobs); l != 0 {
		t.Fatalf("expected to dequeue 0 jobs, got %d", l)
	}

	now = now.Add(40 * time.Second)
	ctx = internal.SetMockTime(ctx, now)
	deqRes = mustDequeueJobs(ctx, t, be, 1, &resource.JobListParams{
		Queues: []string{"cool"},
	})
	if deqRes == nil {
		t.Fatal("expected dequeue result not to be nil")
	}
	if l := len(deqRes.Jobs); l != 1 {
		t.Fatalf("expected to dequeue 1 job, got %d", l)
	}
	// fail for the final time (handler will set dead status in practice)
	now = now.Add(1 * time.Second)
	ctx = internal.SetMockTime(ctx, now)
	mustAckJobs(ctx, t, be, []*resource.Ack{
		{
			JobID:  deqRes.Jobs[0].ID,
			Status: resource.NewStatus(resource.StatusDead),
		},
	})

	// should not retry this again, even after a long time
	now = now.Add(1 * time.Second)
	ctx = internal.SetMockTime(ctx, now)
	deqRes = mustDequeueJobs(ctx, t, be, 1, &resource.JobListParams{
		Queues: []string{"cool"},
	})
	if deqRes == nil {
		t.Fatal("expected dequeue result not to be nil")
	}
	if l := len(deqRes.Jobs); l != 0 {
		t.Fatalf("expected to dequeue 0 jobs, got %d", l)
	}

	now = now.Add(60 * time.Minute)
	ctx = internal.SetMockTime(ctx, now)
	deqRes = mustDequeueJobs(ctx, t, be, 1, &resource.JobListParams{
		Queues: []string{"cool"},
	})
	if deqRes == nil {
		t.Fatal("expected dequeue result not to be nil")
	}
	if l := len(deqRes.Jobs); l != 0 {
		t.Fatalf("expected to dequeue 0 jobs, got %d", l)
	}

	resJobs := mustListJobs(ctx, t, be, 1, &resource.JobListParams{
		Queues: []string{"cool"},
	})
	if l := len(resJobs.Jobs); l != 1 {
		t.Fatalf("expected 1 job, got %d", l)
	}
	status := resJobs.Jobs[0].Status
	if *status != resource.StatusDead {
		t.Fatalf("expected status \"dead\", got %q", status.String())
	}
}

func testDequeueWhileRunning(ctx context.Context, t *testing.T, tc *backendTestContext) {
	be := tc.cfg.Backend
	if tc.cfg.Type == "memory" {
		t.Skip("skipping as memory backend doesn't currently support attempts")
	}
	ctx = mustReset(ctx, t, be)

	now := basictime
	ctx = internal.SetMockTime(ctx, now)
	q := &resource.Queue{
		Name:           "cool",
		Duration:       resource.Duration(1 * time.Minute),
		Retries:        2,
		BackoffInitial: resource.Duration(50 * time.Second),
		BackoffFactor:  2.0,
		BackoffMax:     resource.Duration(10 * time.Minute),
	}
	mustSaveQueue(ctx, t, be, q)

	jobs := getBasicJobs()
	jobs.Jobs = []*resource.Job{jobs.Jobs[0]}
	res := mustEnqueueJobs(ctx, t, be, jobs)
	if len(res.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(res.Jobs))
	}

	now = now.Add(1 * time.Second)
	ctx = internal.SetMockTime(ctx, now)
	deqRes := mustDequeueJobs(ctx, t, be, 1, &resource.JobListParams{
		Queues: []string{"cool"},
	})
	if deqRes == nil {
		t.Fatal("expected dequeue result not to be nil")
	}
	if l := len(deqRes.Jobs); l != 1 {
		t.Fatalf("expected to dequeue 1 job, got %d", l)
	}

	// shouldn't get anything back yet
	now = now.Add(51 * time.Second)
	ctx = internal.SetMockTime(ctx, now)
	deqRes = mustDequeueJobs(ctx, t, be, 1, &resource.JobListParams{
		Queues: []string{"cool"},
	})
	if deqRes == nil {
		t.Fatal("expected dequeue result not to be nil")
	}
	if l := len(deqRes.Jobs); l != 0 {
		t.Fatalf("expected to dequeue 0 jobs, got %d", l)
	}

	// duration + start delay of 1 second + already waited 50 secs
	now = now.Add(10 * time.Second)
	ctx = internal.SetMockTime(ctx, now)
	deqRes = mustDequeueJobs(ctx, t, be, 1, &resource.JobListParams{
		Queues: []string{"cool"},
	})
	if deqRes == nil {
		t.Fatal("expected dequeue result not to be nil")
	}
	if l := len(deqRes.Jobs); l != 1 {
		t.Fatalf("expected to dequeue 1 job, got %d", l)
	}

	// we should get it again if it hasn't been acked by the end of its
	// duration.
	now = now.Add(61 * time.Second)
	ctx = internal.SetMockTime(ctx, now)
	deqRes = mustDequeueJobs(ctx, t, be, 1, &resource.JobListParams{
		Queues: []string{"cool"},
	})
	if deqRes == nil {
		t.Fatal("expected dequeue result not to be nil")
	}
	if l := len(deqRes.Jobs); l != 1 {
		t.Fatalf("expected to dequeue 1 job, got %d", l)
	}
}

func getBasicQueue() *resource.Queue {
	return &resource.Queue{
		Name: "cool",
		// Version:     resource.NewVersion(1),
		Retries: 3,
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

// runMiddleware tries to wrap a call to a backend with its middleware. It
// shouldn't be relied upon for transaction-type middleware.
func runMiddleware(ctx context.Context, t testing.TB, be backend.Interface) (context.Context, func(t testing.TB, err error)) {
	t.Helper()
	origCtx := ctx

	ctxC := make(chan context.Context)
	done := make(chan error)
	mwDone := make(chan struct{})

	if mwer, ok := be.(backend.MiddlewareProvider); ok {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctxC <- r.Context()

			if err := <-done; err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		})

		go func() {
			mockH := mwer.Middleware()(h)
			mockH.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
			mwDone <- struct{}{}
		}()
		// mw := mwer.Middleware()

		ctx = <-ctxC
	} else {
		go func() {
			<-done
			mwDone <- struct{}{}
		}()
	}

	ctx = internal.SetTimeProvider(ctx, internal.GetTimeProvider(origCtx))
	// ctx = internal.SetTicker(ctx, internal.GetTicker(origCtx))
	return ctx, func(t testing.TB, err error) {
		t.Helper()
		done <- err
		<-mwDone
	}
}

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
	ctx, done := runMiddleware(ctx, t, be)

	t.Logf("SaveQueue(%+v)", readable(q))
	res, err := be.SaveQueue(ctx, q)
	if err != nil {
		t.Logf("-> err: %v", err)
		done(t, err)
		t.Fatal(err)
	}
	t.Logf("-> %s", readable(res))
	done(t, nil)
	return res
}

func mustGetQueue(ctx context.Context, t testing.TB, be backend.Interface, q string, opts *resource.GetByIDOpts) *resource.Queue {
	t.Helper()
	ctx, done := runMiddleware(ctx, t, be)

	t.Logf("GetQueue(%+v)", readable(q))
	res, err := be.GetQueue(ctx, q, opts)
	if err != nil {
		t.Logf("-> err: %v", err)
		done(t, err)
		t.Fatal(err)
	}
	t.Logf("-> %s", readable(res))
	done(t, nil)
	return res
}

func checkQueueNotFound(ctx context.Context, t testing.TB, be backend.Interface, q string) {
	t.Helper()
	ctx, done := runMiddleware(ctx, t, be)

	t.Logf("GetQueue(%+v)", readable(q))
	_, err := be.GetQueue(ctx, q, nil)
	if !errors.Is(err, backend.ErrNotFound) {
		done(t, err)
		t.Fatal("expected not found error, got", err)
	}
	t.Logf("-> err (expected): backend.ErrNotFound")
	done(t, nil)
}

func mustDeleteQueue(ctx context.Context, t testing.TB, be backend.Interface, q string) {
	t.Helper()
	ctx, done := runMiddleware(ctx, t, be)

	t.Logf("DeleteQueues(%+v)", readable(q))
	err := be.DeleteQueues(ctx, []string{q})
	if err != nil {
		t.Logf("-> err: %v", err)
		done(t, err)
		t.Fatal(err)
	}
	t.Logf("-> success")
	done(t, nil)
}

func mustEnqueueJobs(ctx context.Context, t testing.TB, be backend.Interface, jobs *resource.Jobs) *resource.Jobs {
	t.Helper()
	ctx, done := runMiddleware(ctx, t, be)

	t.Logf("EnqueueJobs(%+v)", readable(jobs.Jobs))
	res, err := be.EnqueueJobs(ctx, jobs)
	if err != nil {
		t.Logf("-> Error: %v", err)
		done(t, err)
		t.Fatal(err)
	}
	t.Logf("-> %+v", readable(res.Jobs))

	done(t, nil)
	return res
}

func mustDequeueJobs(ctx context.Context, t testing.TB, be backend.Interface, limit int, opts *resource.JobListParams) *resource.Jobs {
	t.Helper()
	ctx, done := runMiddleware(ctx, t, be)

	t.Logf("DequeueJobsOpts(%d, %+v)", limit, readable(opts))
	res, err := be.DequeueJobs(ctx, limit, &resource.JobListParams{
		Queues: []string{"cool"},
	})
	if err != nil {
		t.Logf("-> Error: %v", err)
		done(t, err)
		t.Fatal(err)
	}
	t.Logf("-> %+v", readable(res.Jobs))

	done(t, nil)
	return res
}

func mustAckJobs(ctx context.Context, t testing.TB, be backend.Interface, acks []*resource.Ack) {
	t.Helper()
	ctx, done := runMiddleware(ctx, t, be)

	t.Logf("AckJobs(%+v)", readable(acks))
	err := be.AckJobs(ctx, &resource.Acks{Acks: acks})
	if err != nil {
		t.Logf("-> Error: %v", err)
		done(t, err)
		t.Fatal(err)
	}
	t.Logf("-> OK")
	done(t, nil)
}

func mustListJobs(ctx context.Context, t testing.TB, be backend.Interface, limit int, opts *resource.JobListParams) *resource.Jobs {
	t.Helper()
	ctx, done := runMiddleware(ctx, t, be)

	t.Logf("ListJobs(%d, %+v)", limit, readable(opts))
	res, err := be.ListJobs(ctx, limit, opts)
	if err != nil {
		t.Logf("-> Error: %v", err)
		done(t, err)
		t.Fatal(err)
	}
	t.Logf("-> %+v", readable(res.Jobs))

	done(t, nil)
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

	if q.Name == "" {
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
	} // else if v.Raw() == 1 {
	// }
	return !t.Failed()
}

func checkJob(t testing.TB, jb *resource.Job) bool {
	t.Helper()

	if jb == nil {
		t.Fatal("expected job not to be nil")
		return false
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
	if st := jb.Status; *st != expect {
		t.Errorf("expected job status %s, got %s", resource.NewStatus(expect), st.String())
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
		panic(err)
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
