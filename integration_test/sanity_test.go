package integration

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/qri-io/jsonschema"

	"github.com/jeffrom/job-manager/jobclient"
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

	ctx *sanityContext
}

func (tc *sanityTestCase) wrap(ctx context.Context, fn func(ctx context.Context, t *testing.T, tc *sanityTestCase)) func(t *testing.T) {
	return func(t *testing.T) {
		fn(ctx, t, tc)
	}
}

func (tc *sanityTestCase) saveQueue(ctx context.Context, t testing.TB, name string, opts jobclient.SaveQueueOpts) *jobv1.Queue {
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

func (tc *sanityTestCase) dequeueJobs(ctx context.Context, t testing.TB, num int, name string, selectors ...string) *jobv1.Jobs {
	t.Helper()
	jobs, err := tc.ctx.client.DequeueJobs(ctx, num, name, selectors...)
	if err != nil {
		t.Fatal(err)
	}
	return jobs
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

			// handle multiple jobs
			if !t.Run("handle-multiple", tc.wrap(ctx, testHandleMultipleJobs)) {
				return
			}

			// validate
			if !t.Run("validate-args", tc.wrap(ctx, testValidateArgs)) {
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
	tc.ackJobOpts(ctx, t, id, jobv1.StatusComplete, jobclient.AckJobOpts{
		Data: jobArg.Args[0].GetStringValue(),
	})

	jobData := tc.getJob(ctx, t, id)
	checkJob(t, jobData)

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

// TODO move this to its own test
func testValidateArgs(ctx context.Context, t *testing.T, tc *sanityTestCase) {
	q := tc.saveQueue(ctx, t, "validat0r", jobclient.SaveQueueOpts{
		Schema: testenv.ReadFile(t, "testdata/schema/basic.jsonschema"),
	})
	if q == nil {
		t.Fatal("no queue was saved")
	}

	vtcs := []struct {
		name    string
		queue   string
		args    []interface{}
		keyErrs []jsonschema.KeyError
		errs    []*resource.ValidationError
	}{
		{
			name:  "basic",
			queue: "validat0r",
			args:  jobArgs(-1),
			errs: []*resource.ValidationError{
				{
					Path:  "/0",
					Value: -1,
				},
				{
					Path: "/",
				},
			},
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
			errs: []*resource.ValidationError{
				{
					Path: "/",
				},
			},
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
			errs: []*resource.ValidationError{
				{
					Path:  "/",
					Value: jobArgs("nice"),
				},
			},
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
			errs: []*resource.ValidationError{
				{
					Path:  "/1",
					Value: -1,
				},
			},
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
			errs: []*resource.ValidationError{
				{
					Path:  "/0",
					Value: true,
				},
				{
					Path:  "/1",
					Value: -1,
				},
			},
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
			verrs := getValidationErrors(ctx, t, tc, vtc.queue, vtc.args...)
			if len(verrs) != len(vtc.errs) {
				t.Fatalf("expected %d ValidationErrors, got %d (%+v)", len(vtc.errs), len(verrs), verrs)
			}
			for i, expectErr := range vtc.errs {
				verr := verrs[i]
				if path := expectErr.Path; path != "" {
					checkArgsSchema(t, verr, path)
				} else if len(vtc.args) == 0 && verr.Path != "/" {
					t.Errorf("#%d: expected path to be %q, was %q", i, "/", verr.Path)
				}

				if ival := expectErr.Value; ival != nil {
					switch val := ival.(type) {
					case int:
						checkSchemaNumber(t, verr.Value, float64(val))
					case float64:
						checkSchemaNumber(t, verr.Value, val)
					case string:
						checkSchemaString(t, verr.Value, val)
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
	tc.ackJob(ctx, t, id, jobv1.StatusComplete)
}

func testClaims(ctx context.Context, t *testing.T, tc *sanityTestCase) {
	q := tc.saveQueue(ctx, t, "claimz", jobclient.SaveQueueOpts{
		ClaimDuration: 1 * time.Second,
	})
	if q == nil {
		t.Fatal("no queue was saved")
	}

}

func getValidationErrors(ctx context.Context, t testing.TB, tc *sanityTestCase, queue string, args ...interface{}) []*resource.ValidationError {
	t.Helper()
	id, err := tc.ctx.client.EnqueueJob(ctx, queue, args...)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if id != "" {
		t.Fatal("expected empty id, got", id)
	}

	rerr := &resource.Error{}
	if !errors.As(err, &rerr) {
		t.Fatalf("expected error type %T, got %#v", rerr, err)
	}

	b, err := json.Marshal(rerr.Invalid)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("validation errors: %s", string(b))
	return rerr.Invalid
}

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

func jobArgs(args ...interface{}) []interface{} { return args }
