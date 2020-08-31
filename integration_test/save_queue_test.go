package integration

import (
	"context"
	"testing"

	"github.com/jeffrom/job-manager/jobclient"
	"github.com/jeffrom/job-manager/pkg/job"
	"github.com/jeffrom/job-manager/pkg/testenv"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func TestIntegrationSaveQueue(t *testing.T) {
	tcs := []struct {
		name       string
		argSchema  string
		dataSchema string
		resSchema  string
	}{
		{
			name:      "arg-empty-object",
			argSchema: "{}",
		},
		{
			name:      "arg-type-missing",
			argSchema: `{"minItems": 2}`,
		},
		{
			name:      "arg-type-object",
			argSchema: `{"type": "object"}`,
		},
		{
			name:      "arg-type-invalid",
			argSchema: `{"type": "arrayy"}`,
		},
		{
			name:      "arg-extra",
			argSchema: `{"type": "array", "zorp": true}`,
		},
		{
			name:      "arg-minItems-string",
			argSchema: `{"type": "array", "minItems": "boop"}`,
		},
		{
			name:      "arg-maxItems-string",
			argSchema: `{"type": "array", "maxItems": "boop"}`,
		},
		{
			name:      "arg-items-string",
			argSchema: `{"type": "array", "items": "boop"}`,
		},
		{
			name:      "arg-items-object",
			argSchema: `{"type": "array", "items": {}}`,
		},
		{
			name:       "data-not-type-object",
			dataSchema: `{"type": "array"}`,
		},
		{
			name:       "data-properties-array",
			dataSchema: `{"type": "object", "properties": []}`,
		},
		{
			name:      "result-not-type-object",
			resSchema: `{"type": "array"}`,
		},
		{
			name:      "result-properties-array",
			resSchema: `{"type": "object", "properties": []}`,
		},
	}

	srv := testenv.NewTestControllerServer(t, middleware.NewConfig())
	c := testenv.NewTestClient(t, srv)
	srv.Start()
	defer srv.Close()
	ctx := context.Background()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			checkSaveInvalidQueue(ctx, t, c, tc.name, jobclient.SaveQueueOptions{
				ArgSchema:    []byte(tc.argSchema),
				DataSchema:   []byte(tc.dataSchema),
				ResultSchema: []byte(tc.resSchema),
			})
		})
	}
}

func checkSaveQueue(ctx context.Context, t testing.TB, c jobclient.Interface, name string, opts jobclient.SaveQueueOptions) *job.Queue {
	t.Helper()
	q, err := c.SaveQueue(ctx, name, opts)
	if err != nil {
		t.Fatal(err)
	}
	return q
}

func checkSaveInvalidQueue(ctx context.Context, t testing.TB, c jobclient.Interface, name string, opts jobclient.SaveQueueOptions) {
	t.Helper()
	_, err := c.SaveQueue(ctx, name, opts)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}
