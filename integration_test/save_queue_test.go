package integration

import (
	"context"
	"testing"

	"github.com/jeffrom/job-manager/mjob/client"
	bememory "github.com/jeffrom/job-manager/pkg/backend/mem"
	"github.com/jeffrom/job-manager/pkg/testenv"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func TestIntegrationSaveQueue(t *testing.T) {
	tcs := []struct {
		name   string
		schema string
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

	srv := testenv.NewTestControllerServer(t, middleware.NewConfig(), bememory.New())
	c := testenv.NewTestClient(t, srv)
	srv.Start()
	defer srv.Close()
	ctx := context.Background()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			checkSaveInvalidQueue(ctx, t, c, tc.name, client.SaveQueueOpts{
				Schema: []byte(tc.schema),
			})
		})
	}
}

func checkSaveInvalidQueue(ctx context.Context, t testing.TB, c client.Interface, name string, opts client.SaveQueueOpts) {
	t.Helper()
	_, err := c.SaveQueue(ctx, name, opts)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}
