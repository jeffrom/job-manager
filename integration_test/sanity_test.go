package integration

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/jeffrom/job-manager/jobclient"
	"github.com/jeffrom/job-manager/pkg/job"
	"github.com/jeffrom/job-manager/pkg/schema"
	"github.com/jeffrom/job-manager/pkg/testenv"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
	"github.com/qri-io/jsonschema"
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

func (tc *sanityTestCase) saveQueue(ctx context.Context, t testing.TB, name string, opts jobclient.SaveQueueOpts) *job.Queue {
	t.Helper()
	q, err := tc.ctx.client.SaveQueue(ctx, name, opts)
	if err != nil {
		t.Fatal(err)
	}
	return q
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
			if !t.Run("validate-args", tc.wrap(ctx, testValidateArgs)) {
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
	expectErr := &jobclient.NotFoundError{}
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

	if len(seen) != n {
		t.Errorf("expected %d dequeued job ids, got %d (%v)", n, len(seen), seen)
	}

	jobs = tc.dequeueJobs(ctx, t, n, "cool")
	if len(jobs.Jobs) != 0 {
		t.Fatalf("expected 0 jobs, got %d", len(jobs.Jobs))
	}

	for _, id := range ids {
		tc.ackJob(ctx, t, id, job.StatusFailed)
	}

	jobs = tc.dequeueJobs(ctx, t, n, "cool")
	if len(jobs.Jobs) != n {
		t.Fatalf("expected %d jobs, got %d", n, len(jobs.Jobs))
	}
	checkJob(t, jobs.Jobs[0])
	checkJob(t, jobs.Jobs[1])
	checkJob(t, jobs.Jobs[2])

	for _, id := range ids {
		tc.ackJob(ctx, t, id, job.StatusComplete)
	}

	jobs = tc.dequeueJobs(ctx, t, n, "cool")
	if len(jobs.Jobs) != 0 {
		t.Fatalf("expected 0 jobs, got %d", len(jobs.Jobs))
	}
}

// TODO move this to its own test
func testValidateArgs(ctx context.Context, t *testing.T, tc *sanityTestCase) {
	q := tc.saveQueue(ctx, t, "validat0r", jobclient.SaveQueueOpts{
		ArgSchema: testenv.ReadFile(t, "testdata/schema/basic.jsonschema"),
	})
	if q == nil {
		t.Fatal("no queue was saved")
	}

	vtcs := []struct {
		name    string
		queue   string
		args    []interface{}
		keyErrs []jsonschema.KeyError
	}{
		{
			name:  "basic",
			queue: "validat0r",
			args:  jobArgs(-1),
			keyErrs: []jsonschema.KeyError{
				{
					PropertyPath: "/0",
					InvalidValue: -1,
				},
				{
					PropertyPath: "/",
				},
			},
		},
		{
			name:  "noargs",
			queue: "validat0r",
			args:  jobArgs(),
			keyErrs: []jsonschema.KeyError{
				{
					PropertyPath: "/",
				},
			},
		},
		{
			name:  "incomplete",
			queue: "validat0r",
			args:  jobArgs("nice"),
			keyErrs: []jsonschema.KeyError{
				{
					PropertyPath: "/",
					// TODO implement interface array check
					InvalidValue: jobArgs("nice"),
				},
			},
		},
		{
			name:  "wrong-type",
			queue: "validat0r",
			args:  jobArgs("nice", -1),
			keyErrs: []jsonschema.KeyError{
				{
					PropertyPath: "/1",
					InvalidValue: -1,
				},
			},
		},
		{
			name:  "wrong-types",
			queue: "validat0r",
			args:  jobArgs(true, -1),
			keyErrs: []jsonschema.KeyError{
				{
					PropertyPath: "/0",
					InvalidValue: true,
				},
				{
					PropertyPath: "/1",
					InvalidValue: -1,
				},
			},
		},
	}
	for _, vtc := range vtcs {
		t.Run(vtc.name, func(t *testing.T) {
			t.Log(vtc.queue, vtc.args)
			keyErrs := getKeyErrors(ctx, t, tc, vtc.queue, vtc.args...)
			if len(keyErrs) != len(vtc.keyErrs) {
				t.Fatalf("expected %d jsonschema KeyError, got %d (%s)", len(vtc.keyErrs), len(keyErrs), keyErrs)
			}
			for i, expectErr := range vtc.keyErrs {
				keyErr := keyErrs[i]
				if path := expectErr.PropertyPath; path != "" {
					checkArgsSchema(t, keyErr, path)
				} else if len(vtc.args) == 0 && keyErr.PropertyPath != "/" {
					t.Errorf("#%d: expected path to be %q, was %q", i, "/", keyErr.PropertyPath)
				}

				if ival := expectErr.InvalidValue; ival != nil {
					switch val := ival.(type) {
					case int:
						checkSchemaNumber(t, keyErr.InvalidValue, float64(val))
					case float64:
						checkSchemaNumber(t, keyErr.InvalidValue, val)
					case string:
						checkSchemaString(t, keyErr.InvalidValue, val)
					}
				}
			}

			// for i, iArg := range vtc.args {
			// 	keyErr := keyErrs[i]
			// 	switch vArg := iArg.(type) {
			// 	case float64:
			// 		checkSchemaNumber(t, keyErr.InvalidValue, vArg)
			// 	}
			// }
		})
	}

	id := tc.enqueueJob(ctx, t, "validat0r", "nice", true)
	tc.ackJob(ctx, t, id, job.StatusComplete)
}

func getKeyErrors(ctx context.Context, t testing.TB, tc *sanityTestCase, queue string, args ...interface{}) []jsonschema.KeyError {
	t.Helper()
	id, err := tc.ctx.client.EnqueueJob(ctx, queue, args...)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if id != "" {
		t.Fatal("expected empty id, got", id)
	}

	verr := &schema.ValidationError{}
	if !errors.As(err, &verr) {
		t.Fatalf("expected error type %T, got %#v", verr, err)
	}

	return verr.KeyErrors()
}

func checkArgsSchema(t testing.TB, keyErr jsonschema.KeyError, expectPath string) {
	t.Helper()
	if keyErr.PropertyPath != expectPath {
		t.Errorf("expected path %q, got %q", expectPath, keyErr.PropertyPath)
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

func jobArgs(args ...interface{}) []interface{} { return args }
