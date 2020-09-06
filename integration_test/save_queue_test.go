package integration

import (
	"context"
	"testing"

	"github.com/jeffrom/job-manager/jobclient"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
	"github.com/jeffrom/job-manager/pkg/testenv"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func TestIntegrationSaveQueue(t *testing.T) {
	tcs := []struct {
		name       string
		schema     string
		argSchema  string
		dataSchema string
		resSchema  string
	}{
		{
			name:   "arg-empty-object",
			schema: `{"args": {}}`,
		},
		{
			name:   "arg-type-missing",
			schema: `{"args": {"minItems": 2}}`,
		},
		{
			name:   "arg-type-object",
			schema: `{"args": {"type": "object"}}`,
		},
		{
			name:   "arg-type-invalid",
			schema: `{"args": {"type": "arrayy"}}`,
		},
		{
			name:   "arg-extra",
			schema: `{"args": {"type": "array", "zorp": true}}`,
		},
		{
			name:   "arg-minItems-string",
			schema: `{"args": {"type": "array", "minItems": "boop"}}`,
		},
		{
			name:   "arg-maxItems-string",
			schema: `{"args": {"type": "array", "maxItems": "boop"}}`,
		},
		{
			name:   "arg-items-string",
			schema: `{"args": {"type": "array", "items": "boop"}}`,
		},
		{
			name:   "arg-items-object",
			schema: `{"args": {"type": "array", "items": {}}}`,
		},
		{
			name:   "data-not-type-object",
			schema: `{"data": {"type": "array"}}`,
		},
		{
			name:   "data-properties-array",
			schema: `{"data": {"type": "object", "properties": []}}`,
		},
		{
			name:   "result-not-type-object",
			schema: `{"result": {"type": "array"}}`,
		},
		{
			name:   "result-properties-array",
			schema: `{"result": {"type": "object", "properties": []}}`,
		},
	}

	srv := testenv.NewTestControllerServer(t, middleware.NewConfig())
	c := testenv.NewTestClient(t, srv)
	srv.Start()
	defer srv.Close()
	ctx := context.Background()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			checkSaveInvalidQueue(ctx, t, c, tc.name, jobclient.SaveQueueOpts{
				Schema: []byte(tc.schema),
				// ArgSchema:    []byte(tc.argSchema),
				// DataSchema:   []byte(tc.dataSchema),
				// ResultSchema: []byte(tc.resSchema),
			})
		})
	}
}

func checkSaveQueue(ctx context.Context, t testing.TB, c jobclient.Interface, name string, opts jobclient.SaveQueueOpts) *jobv1.Queue {
	t.Helper()
	q, err := c.SaveQueue(ctx, name, opts)
	if err != nil {
		t.Fatal(err)
	}
	return q
}

func checkSaveInvalidQueue(ctx context.Context, t testing.TB, c jobclient.Interface, name string, opts jobclient.SaveQueueOpts) {
	t.Helper()
	_, err := c.SaveQueue(ctx, name, opts)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}
