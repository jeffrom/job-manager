package integration

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jeffrom/job-manager/jobclient"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/label"
	"github.com/jeffrom/job-manager/pkg/resource"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
	"github.com/jeffrom/job-manager/pkg/testenv"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

type sanityContext struct {
	srv    *httptest.Server
	client jobclient.Interface
}

type sanityTestCase struct {
	name   string
	srvCfg middleware.Config
	now    time.Time

	ctx *sanityContext
}

func (tc *sanityTestCase) wrap(ctx context.Context, fn func(ctx context.Context, t *testing.T, tc *sanityTestCase)) func(t *testing.T) {
	return func(t *testing.T) {
		fn(ctx, t, tc)
	}
}

func (tc *sanityTestCase) setMockTime(ctx context.Context, t testing.TB, ts time.Time) context.Context {
	t.Helper()
	t.Logf("setting mock time: %s (%d)", ts.Format(time.Stamp), ts.Unix())
	tc.now = ts
	return jobclient.SetMockTime(ctx, ts)
}

func (tc *sanityTestCase) incMockTime(ctx context.Context, t testing.TB, dur time.Duration) (context.Context, time.Time) {
	tc.now = tc.now.Add(dur)
	t.Logf("incrementing mock time to: %s (%d) (+%s)", tc.now.Format(time.Stamp), tc.now.Unix(), dur.String())
	ctx = jobclient.SetMockTime(ctx, tc.now)
	return ctx, tc.now
}

func (tc *sanityTestCase) saveQueue(ctx context.Context, t testing.TB, name string, opts jobclient.SaveQueueOpts) *jobv1.Queue {
	t.Helper()
	t.Logf("SaveQueue(%q)", name)
	q, err := tc.ctx.client.SaveQueue(ctx, name, opts)
	if err != nil {
		t.Logf("-> Error: %v", err)
		t.Fatal(err)
	}
	t.Logf("-> %s", q.String())
	return q
}

func (tc *sanityTestCase) enqueueJobOpts(ctx context.Context, t testing.TB, name string, opts jobclient.EnqueueOpts, args ...interface{}) string {
	t.Helper()
	t.Logf("EnqueueJobOpts(%q, %+v, %+v)", name, opts, args)
	id, err := tc.ctx.client.EnqueueJobOpts(ctx, name, opts, args...)
	if err != nil {
		t.Logf("-> Error: %v", err)
		t.Fatal(err)
	}
	t.Logf("-> %s", id)
	return id
}

func (tc *sanityTestCase) enqueueJob(ctx context.Context, t testing.TB, name string, args ...interface{}) string {
	return tc.enqueueJobOpts(ctx, t, name, jobclient.EnqueueOpts{}, args...)
}

func (tc *sanityTestCase) dequeueJobsOpts(ctx context.Context, t testing.TB, num int, opts jobclient.DequeueOpts) *jobv1.Jobs {
	t.Helper()
	jobs, err := tc.ctx.client.DequeueJobsOpts(ctx, num, opts)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("dequeued %d %s", len(jobs.Jobs), jobs.String())
	return jobs
}

func (tc *sanityTestCase) dequeueJobs(ctx context.Context, t testing.TB, num int, name string) *jobv1.Jobs {
	t.Helper()
	return tc.dequeueJobsOpts(ctx, t, num, jobclient.DequeueOpts{Queues: []string{name}})
}

func (tc *sanityTestCase) ackJob(ctx context.Context, t testing.TB, id string, status jobv1.Status) {
	t.Helper()
	if err := tc.ctx.client.AckJob(ctx, id, status); err != nil {
		t.Fatal(err)
	}
}

func (tc *sanityTestCase) ackJobOpts(ctx context.Context, t testing.TB, id string, status jobv1.Status, opts jobclient.AckJobOpts) {
	t.Helper()
	if err := tc.ctx.client.AckJobOpts(ctx, id, status, opts); err != nil {
		t.Fatal(err)
	}
	t.Logf("ack %s: %s", id, status.String())
}

func (tc *sanityTestCase) getJob(ctx context.Context, t testing.TB, id string) *jobv1.Job {
	t.Helper()
	jobData, err := tc.ctx.client.GetJob(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	return jobData
}

// TestIntegrationSanity goes through the basic operations (soon with a variety
// of configs)
func TestIntegrationSanity(t *testing.T) {
	tcs := []sanityTestCase{
		{
			name:   "default",
			srvCfg: middleware.NewConfig(),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			srv := testenv.NewTestControllerServer(t, tc.srvCfg, backend.NewMemory())
			c := testenv.NewTestClient(t, srv)
			tc.ctx = &sanityContext{srv: srv, client: c}
			srv.Start()
			defer srv.Close()

			ctx := context.Background()
			t.Run("ping", tc.wrap(ctx, testPing))

			// enqueue a job, dequeue it
			if !t.Run("single-job", tc.wrap(ctx, testSingleJob)) {
				return
			}

			// handle multiple jobs
			if !t.Run("handle-multiple", tc.wrap(ctx, testHandleMultipleJobs)) {
				return
			}

			// unique jobs
			// unique jobs (jsonschema)

			// checkins

			// claim windows
			if !t.Run("claims", tc.wrap(ctx, testClaims)) {
				return
			}
		})
	}
}

func testPing(ctx context.Context, t *testing.T, tc *sanityTestCase) {
	c := tc.ctx.client
	if err := c.Ping(ctx); err != nil {
		t.Fatal(err)
	}
}

func testSingleJob(ctx context.Context, t *testing.T, tc *sanityTestCase) {
	// NOTE this block of tests shares state
	t.Run("enqueue-no-queue", tc.wrap(ctx, testEnqueueNoQueue))
	t.Run("create-queue", tc.wrap(ctx, testCreateQueue))
	t.Run("dequeue-no-jobs", tc.wrap(ctx, testDequeueEmpty))
	t.Run("enqueue", tc.wrap(ctx, testEnqueue))
	t.Run("dequeue", tc.wrap(ctx, testDequeue))
	t.Run("dequeue-empty", tc.wrap(ctx, testDequeueEmpty))
}

func testEnqueueNoQueue(ctx context.Context, t *testing.T, tc *sanityTestCase) {
	c := tc.ctx.client
	expectErr := &resource.Error{
		Status:     404,
		Kind:       "not_found",
		Resource:   "queue",
		ResourceID: "cool",
	}
	id, err := c.EnqueueJob(ctx, "cool", "nice")
	if !errors.Is(err, expectErr) {
		// fmt.Printf("e: %#v\n", expectErr)
		// fmt.Printf("g: %#v\n", err)
		t.Errorf("expected error %T, got %#v", expectErr, err)
	}
	if id != "" {
		t.Error("expected empty id, got", id)
	}
}

func testCreateQueue(ctx context.Context, t *testing.T, tc *sanityTestCase) {
	c := tc.ctx.client
	expectID := "cool"
	q, err := c.SaveQueue(ctx, expectID, jobclient.SaveQueueOpts{})
	if err != nil {
		t.Fatal(err)
	}

	if q == nil {
		t.Fatal("queue result was nil")
	}
	if q.Id != expectID {
		t.Errorf("expected queue name %q, got %q", expectID, q.Id)
	}
	var defaultConcurrency int32 = 10
	if q.Concurrency != defaultConcurrency {
		t.Errorf("expected default concurrency %d, got %d", defaultConcurrency, q.Concurrency)
	}
	var defaultMaxRetries int32 = 10
	if q.Retries != defaultMaxRetries {
		t.Errorf("expected default max retries %d, got %d", defaultMaxRetries, q.Retries)
	}

	if q.V != 1 {
		t.Fatalf("expected queue version to be %d, was %d", 1, q.V)
	}
}

var basicTime = time.Date(2020, 1, 1, 13, 0, 0, 0, time.UTC)

func testEnqueue(ctx context.Context, t *testing.T, tc *sanityTestCase) {
	c := tc.ctx.client
	now := basicTime
	ctx = tc.setMockTime(ctx, t, now)
	id, err := c.EnqueueJob(ctx, "cool", "nice")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("got job id: %q", id)
	if id == "" {
		t.Fatal("id was empty")
	}
}

func testDequeue(ctx context.Context, t *testing.T, tc *sanityTestCase) {
	expectJobName := "cool"
	ctx, now := tc.incMockTime(ctx, t, 1*time.Second)
	jobs := tc.dequeueJobs(ctx, t, 1, expectJobName)
	if len(jobs.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs.Jobs))
	}

	jobArg := jobs.Jobs[0]
	// t.Logf("job: %+v")

	if jobArg.Name != expectJobName {
		t.Errorf("expected job name %q, got %q", expectJobName, jobArg.Name)
	}

	if jobArg.V != 2 {
		t.Errorf("expected job version to be v2, was %d", jobArg.V)
	}

	if len(jobArg.Args) != 1 {
		t.Errorf("expected job args length %d, got %d", 1, len(jobArg.Args))
	} else {
		arg := jobArg.Args[0].GetStringValue()
		expectArg := "nice"
		if arg != expectArg {
			t.Errorf("expected job arg to be %q, was %q", expectArg, arg)
		}
	}

	checkJob(t, jobArg)
	if enqueuedAt := jobArg.EnqueuedAt.AsTime(); !enqueuedAt.Equal(basicTime) {
		t.Errorf("expected enqueued_at to be %q, was %q", basicTime, enqueuedAt)
	}

	res := jobArg.Results
	if len(res) != 1 {
		t.Fatal("expected 1 result, got", len(res))
	}

	resArg := res[0]
	if startedAt := resArg.StartedAt.AsTime(); !startedAt.Equal(now) {
		t.Errorf("expected started_at to be %q, was %q", now, startedAt)
	}
	if completedAt := resArg.CompletedAt.AsTime(); !completedAt.IsZero() {
		t.Error("expected completed_at to be zero")
	}

	id := jobArg.Id
	tc.ackJobOpts(ctx, t, id, jobv1.StatusComplete, jobclient.AckJobOpts{
		Data: jobArg.Args[0].GetStringValue(),
	})

	jobData := tc.getJob(ctx, t, id)
	checkJob(t, jobData)
	if !jobData.EnqueuedAt.AsTime().Equal(basicTime) {
		t.Errorf("expected enqueued_at to be %q, was %q", basicTime, jobData.EnqueuedAt.AsTime())
	}

	resultData := jobData.Results
	if resultData == nil {
		t.Fatal("job result data was nil")
	}
	if len(resultData) != 1 {
		t.Fatalf("expected 1 results, got %d", len(resultData))
	}
	ival := resultData[0].Data
	// ival, ok := resultData["arg"]
	// if !ok {
	// 	t.Fatal("expected result data attribute 'arg'")
	// }
	s, ok := ival.AsInterface().(string)
	if !ok {
		t.Fatalf("expected result data to be type string, was %T", ival)
	}
	if s != "nice" {
		t.Errorf("expected 'arg' to be value 'nice', was %q", s)
	}
}

func testDequeueEmpty(ctx context.Context, t *testing.T, tc *sanityTestCase) {
	c := tc.ctx.client
	jobs, err := c.DequeueJobs(ctx, 1, "cool")
	if err != nil {
		t.Fatal(err)
	}
	if jobs == nil {
		t.Fatal("dequeue result was nil")
	}
	if len(jobs.Jobs) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(jobs.Jobs))
	}
}

func testHandleMultipleJobs(ctx context.Context, t *testing.T, tc *sanityTestCase) {
	n := 3
	ids := make([]string, n)
	for i := 0; i < n; i++ {
		ids[i] = tc.enqueueJob(ctx, t, "cool", i)
	}

	seen := make(map[string]bool)
	jobs := tc.dequeueJobs(ctx, t, 1, "cool")
	if len(jobs.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs.Jobs))
	}
	checkJob(t, jobs.Jobs[0])
	seen[jobs.Jobs[0].Id] = true

	jobs = tc.dequeueJobs(ctx, t, 3, "cool")
	if len(jobs.Jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(jobs.Jobs))
	}
	checkJob(t, jobs.Jobs[0])
	checkJob(t, jobs.Jobs[1])
	seen[jobs.Jobs[0].Id] = true
	seen[jobs.Jobs[1].Id] = true

	if len(seen) != n {
		t.Errorf("expected %d dequeued job ids, got %d (%v)", n, len(seen), seen)
	}

	jobs = tc.dequeueJobs(ctx, t, n, "cool")
	if len(jobs.Jobs) != 0 {
		t.Fatalf("expected 0 jobs, got %d", len(jobs.Jobs))
	}

	for _, id := range ids {
		tc.ackJob(ctx, t, id, jobv1.StatusFailed)
	}

	jobs = tc.dequeueJobs(ctx, t, n, "cool")
	if len(jobs.Jobs) != n {
		t.Fatalf("expected %d jobs, got %d", n, len(jobs.Jobs))
	}
	checkJob(t, jobs.Jobs[0])
	checkJob(t, jobs.Jobs[1])
	checkJob(t, jobs.Jobs[2])

	for _, id := range ids {
		tc.ackJob(ctx, t, id, jobv1.StatusComplete)
	}

	jobs = tc.dequeueJobs(ctx, t, n, "cool")
	if len(jobs.Jobs) != 0 {
		t.Fatalf("expected 0 jobs, got %d", len(jobs.Jobs))
	}
}

// testClaims tests behavior around claim windows. it should work as follows:
//
// - if the claim window has not elapsed, job-manager shouldn't dequeue jobs
// without matching claims.
// - the claim window should be reset when a job fails.
// - if the claim window has elapsed, claims should be ignored.
func testClaims(ctx context.Context, t *testing.T, tc *sanityTestCase) {
	q := tc.saveQueue(ctx, t, "claimz", jobclient.SaveQueueOpts{
		ClaimDuration: 1 * time.Second,
	})
	if q == nil {
		t.Fatal("no queue was saved")
	}
	claims := label.Claims(map[string][]string{
		"coolclaim": []string{"itiscool"},
	})
	mockNow := basicTime
	ctx = tc.setMockTime(ctx, t, mockNow)
	id := tc.enqueueJobOpts(ctx, t, "claimz", jobclient.EnqueueOpts{
		Claims: claims,
	})

	jobs := tc.dequeueJobs(ctx, t, 1, "claimz")
	if len(jobs.Jobs) != 0 {
		t.Fatalf("expected 0 jobs, got %d", len(jobs.Jobs))
	}

	dqClaimOpts := jobclient.DequeueOpts{Claims: claims}
	claimJobs := tc.dequeueJobsOpts(ctx, t, 1, dqClaimOpts)
	if len(claimJobs.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs.Jobs))
	}

	// fail the job
	nextNow := mockNow.Add(1 * time.Second)
	ctx = tc.setMockTime(ctx, t, nextNow)
	tc.ackJobOpts(ctx, t, id, jobv1.StatusFailed, jobclient.AckJobOpts{})

	// try to dequeue again, ensure claim window is reset
	if jobs := tc.dequeueJobs(ctx, t, 1, "claimz"); len(jobs.Jobs) != 0 {
		t.Fatalf("expected 0 jobs, got %d", len(jobs.Jobs))
	}

	// dequeue before claim window has elapsed with claims
	if claimJobs := tc.dequeueJobsOpts(ctx, t, 1, dqClaimOpts); len(claimJobs.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs.Jobs))
	}

	nextNow = nextNow.Add(1 * time.Second)
	ctx = tc.setMockTime(ctx, t, nextNow)
	tc.ackJobOpts(ctx, t, id, jobv1.StatusFailed, jobclient.AckJobOpts{})

	// try to dequeue again, ensure claim window is reset
	if jobs := tc.dequeueJobs(ctx, t, 1, "claimz"); len(jobs.Jobs) != 0 {
		t.Fatalf("expected 0 jobs, got %d", len(jobs.Jobs))
	}

	claimElapsed := nextNow.Add(1 * time.Second)
	ctx = tc.setMockTime(ctx, t, claimElapsed)
	if claimJobs := tc.dequeueJobs(ctx, t, 1, "claimz"); len(claimJobs.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs.Jobs))
	}

	tc.ackJobOpts(ctx, t, id, jobv1.StatusComplete, jobclient.AckJobOpts{
		// Claims: nil,
	})
}

func jobArgs(args ...interface{}) []interface{} { return args }
