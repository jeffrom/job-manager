package integration

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/resource"
	"github.com/jeffrom/job-manager/pkg/testenv"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
	"github.com/qri-io/jsonschema"
)

func TestIntegrationArgsValidate(t *testing.T) {
	tcs := []struct {
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

	srv := testenv.NewTestControllerServer(t, middleware.NewConfig(), backend.NewMemory())
	srv.Start()
	defer srv.Close()

	c := testenv.NewTestClient(t, srv)
	ctx := context.Background()

	q, err := c.SaveQueue(ctx, "validat0r", client.SaveQueueOpts{
		Schema: testenv.ReadFile(t, "testdata/schema/basic.jsonschema"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if q == nil {
		t.Fatal("expected queue to be saved")
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.queue, tc.args)
			verrs := getArgValidationErrors(ctx, t, c, tc.queue, tc.args...)
			if len(verrs) != len(tc.errs) {
				t.Fatalf("expected %d ValidationErrors, got %d (%+v)", len(tc.errs), len(verrs), verrs)
			}
			for i, expectErr := range tc.errs {
				verr := verrs[i]
				if path := expectErr.Path; path != "" {
					checkArgsSchema(t, verr, path)
				} else if len(tc.args) == 0 && verr.Path != "/" {
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
		})
	}

	id, err := c.EnqueueJob(ctx, "validat0r", "nice", true)
	if err != nil {
		t.Fatal(err)
	}

	jobs, err := c.DequeueJobs(ctx, 1, "validat0r")
	if err != nil {
		return
	}

	if len(jobs.Jobs) != 1 {
		t.Fatal("expected to dequeue one job, got", len(jobs.Jobs))
	}

	if err := c.AckJob(ctx, id, resource.StatusComplete); err != nil {
		t.Fatal(err)
	}
}

func getArgValidationErrors(ctx context.Context, t testing.TB, c client.Interface, queue string, args ...interface{}) []*resource.ValidationError {
	t.Helper()

	id, err := c.EnqueueJob(ctx, queue, args...)
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
