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
	"github.com/jeffrom/job-manager/pkg/web"
)

type sanityContext struct {
	srv    *httptest.Server
	client *jobclient.Client
}

type sanityTestCase struct {
	name   string
	srvCfg web.Config

	ctx *sanityContext
}

func (tc *sanityTestCase) wrap(ctx context.Context, fn func(ctx context.Context, t *testing.T, tc *sanityTestCase)) func(t *testing.T) {
	return func(t *testing.T) {
		fn(ctx, t, tc)
	}
}

// TestIntegrationSanity goes through the basic operations (soon with a variety
// of configs)
func TestIntegrationSanity(t *testing.T) {
	tcs := []sanityTestCase{
		{
			name: "default",
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
			t.Run("enqueue-no-queue", tc.wrap(ctx, testEnqueueNoQueue))
			t.Run("create-queue", tc.wrap(ctx, testCreateQueue))
			t.Run("dequeue-no-jobs", tc.wrap(ctx, testDequeueEmpty))
			t.Run("enqueue", tc.wrap(ctx, testEnqueue))
			t.Run("dequeue", tc.wrap(ctx, testDequeue))
			t.Run("dequeue-empty", tc.wrap(ctx, testDequeueEmpty))
		})
	}
}

func testPing(ctx context.Context, t *testing.T, tc *sanityTestCase) {
	c := tc.ctx.client
	if err := c.Ping(ctx); err != nil {
		t.Fatal(err)
	}
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

	job := jobs.Jobs[0]

	if job.Name != expectJobName {
		t.Errorf("expected job name %q, got %q", expectJobName, job.Name)
	}

	if len(job.Args) != 1 {
		t.Errorf("expected job args length %d, got %d", 1, len(job.Args))
	} else {
		arg := job.Args[0].GetStringValue()
		expectArg := "nice"
		if arg != expectArg {
			t.Errorf("expected job arg to be %q, was %q", expectArg, arg)
		}
	}

	sanityCheckJob(t, job)
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

func sanityCheckJob(t testing.TB, jobData *job.Job) {
	if jobData.EnqueuedAt == nil {
		t.Errorf("job.EnqueuedAt was nil")
	} else if !jobData.EnqueuedAt.IsValid() {
		t.Errorf("job.EnqueuedAt is invalid: %v", jobData.EnqueuedAt.CheckValid())
	} else if jobData.EnqueuedAt.AsTime().IsZero() {
		t.Errorf("job.EnqueuedAt is zero")
	}

	if jobData.Status == job.Status_UNKNOWN {
		t.Errorf("job.Status is %s", job.Status_UNKNOWN)
	}
	if jobData.Attempt < 0 {
		t.Errorf("job.Attempt was %d", jobData.Attempt)
	}
	if jobData.MaxRetries < 1 {
		t.Errorf("job.MaxRetries was %d", jobData.MaxRetries)
	}
	if jobData.Duration == nil {
		t.Errorf("job.Duration was nil")
	}
}
