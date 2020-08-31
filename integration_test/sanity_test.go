package integration

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/jeffrom/job-manager/jobclient"
	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/job"
	"github.com/jeffrom/job-manager/pkg/testenv"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

type sanityContext struct {
	srv    *httptest.Server
	client *jobclient.Client
}

type sanityTestCase struct {
	name   string
	srvCfg middleware.Config

	ctx *sanityContext
}

func (tc *sanityTestCase) wrap(ctx context.Context, fn func(ctx context.Context, t *testing.T, tc *sanityTestCase)) func(t *testing.T) {
	return func(t *testing.T) {
		fn(ctx, t, tc)
	}
}

func (tc *sanityTestCase) enqueueJob(ctx context.Context, t testing.TB, name string, args ...interface{}) string {
	t.Helper()
	id, err := tc.ctx.client.EnqueueJob(ctx, name, args...)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func (tc *sanityTestCase) dequeueJobs(ctx context.Context, t testing.TB, num int, name string, selectors ...string) *job.Jobs {
	t.Helper()
	jobs, err := tc.ctx.client.DequeueJobs(ctx, num, name, selectors...)
	if err != nil {
		t.Fatal(err)
	}
	return jobs
}

func (tc *sanityTestCase) ackJob(ctx context.Context, t testing.TB, id string, status job.Status) {
	t.Helper()
	if err := tc.ctx.client.AckJob(ctx, id, status); err != nil {
		t.Fatal(err)
	}
}

func (tc *sanityTestCase) ackJobOpts(ctx context.Context, t testing.TB, id string, status job.Status, opts jobclient.AckJobOpts) {
	t.Helper()
	if err := tc.ctx.client.AckJobOpts(ctx, id, status, opts); err != nil {
		t.Fatal(err)
	}
}

func (tc *sanityTestCase) getJob(ctx context.Context, t testing.TB, id string) *job.Job {
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
			srv := testenv.NewTestControllerServer(t, tc.srvCfg)
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
			if !t.Run("handle-multiple", tc.wrap(ctx, testHandleMultipleJobs)) {
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
	expectErr := &apiv1.NotFoundError{}
	id, err := c.EnqueueJob(ctx, "cool", "nice")
	if !errors.Is(err, expectErr) {
		t.Errorf("expected error %T, got %#v", expectErr, err)
	}
	if id != "" {
		t.Error("expected empty id, got", id)
	}
}

func testCreateQueue(ctx context.Context, t *testing.T, tc *sanityTestCase) {
	c := tc.ctx.client
	expectName := "cool"
	q, err := c.SaveQueue(ctx, expectName, jobclient.SaveQueueOptions{})
	if err != nil {
		t.Fatal(err)
	}

	if q == nil {
		t.Fatal("queue result was nil")
	}
	if q.Name != expectName {
		t.Errorf("expected queue name %q, got %q", expectName, q.Name)
	}
	var defaultConcurrency int32 = 10
	if q.Concurrency != defaultConcurrency {
		t.Errorf("expected default concurrency %d, got %d", defaultConcurrency, q.Concurrency)
	}
	var defaultMaxRetries int32 = 10
	if q.MaxRetries != defaultMaxRetries {
		t.Errorf("expected default max retries %d, got %d", defaultMaxRetries, q.MaxRetries)
	}
}

func testEnqueue(ctx context.Context, t *testing.T, tc *sanityTestCase) {
	c := tc.ctx.client
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
	c := tc.ctx.client
	expectJobName := "cool"
	jobs, err := c.DequeueJobs(ctx, 1, expectJobName)
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs.Jobs))
	}

	jobArg := jobs.Jobs[0]

	if jobArg.Name != expectJobName {
		t.Errorf("expected job name %q, got %q", expectJobName, jobArg.Name)
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

	id := jobArg.Id
	tc.ackJobOpts(ctx, t, id, job.StatusComplete, jobclient.AckJobOpts{
		Data: map[string]interface{}{"arg": jobArg.Args[0].GetStringValue()},
	})

	jobData := tc.getJob(ctx, t, id)
	checkJob(t, jobData)

	resultData := jobData.ResultData
	if resultData == nil {
		t.Fatal("job result data was nil")
	}
	ival, ok := resultData["arg"]
	if !ok {
		t.Fatal("expected result data attribute 'arg'")
	}
	s, ok := ival.AsInterface().(string)
	if !ok {
		t.Fatalf("expected result data attribute 'arg' to be type string, was %T", ival)
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

	if len(seen) != 3 {
		t.Errorf("expected 3 dequeued job ids, got %d (%v)", len(seen), seen)
	}

	jobs = tc.dequeueJobs(ctx, t, 3, "cool")
	if len(jobs.Jobs) != 0 {
		t.Fatalf("expected 0 jobs, got %d", len(jobs.Jobs))
	}

	for _, id := range ids {
		tc.ackJob(ctx, t, id, job.StatusFailed)
	}

	jobs = tc.dequeueJobs(ctx, t, 3, "cool")
	if len(jobs.Jobs) != 3 {
		t.Fatalf("expected 3 jobs, got %d", len(jobs.Jobs))
	}
	checkJob(t, jobs.Jobs[0])
	checkJob(t, jobs.Jobs[1])
	checkJob(t, jobs.Jobs[2])

	for _, id := range ids {
		tc.ackJob(ctx, t, id, job.StatusComplete)
	}

	jobs = tc.dequeueJobs(ctx, t, 3, "cool")
	if len(jobs.Jobs) != 0 {
		t.Fatalf("expected 0 jobs, got %d", len(jobs.Jobs))
	}
}

func checkJob(t testing.TB, jobData *job.Job) {
	if jobData.EnqueuedAt == nil {
		t.Errorf("job.EnqueuedAt was nil")
	} else if !jobData.EnqueuedAt.IsValid() {
		t.Errorf("job.EnqueuedAt is invalid: %v", jobData.EnqueuedAt.CheckValid())
	} else if jobData.EnqueuedAt.AsTime().IsZero() {
		t.Errorf("job.EnqueuedAt is zero")
	}

	if jobData.Status == job.StatusUnknown {
		t.Errorf("job.Status is %s", job.StatusUnknown)
	}
	if jobData.Attempt < 0 {
		t.Errorf("job.Attempt was %d", jobData.Attempt)
	}
	if jobData.MaxRetries < 0 {
		t.Errorf("job.MaxRetries was %d", jobData.MaxRetries)
	}
	if jobData.Duration == nil {
		t.Errorf("job.Duration was nil")
	}
}
