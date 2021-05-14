package mjob

import (
	"context"
	"testing"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/resource"
)

type mockClient struct {
	client.Client
	nextQ  *resource.Queue
	called bool
}

func (mc *mockClient) GetQueue(ctx context.Context, id string) (*resource.Queue, error) {
	mc.called = true
	return mc.nextQ, nil
}

func (mc *mockClient) setQueue(q *resource.Queue) {
	mc.nextQ = q
}

func (mc *mockClient) resetCalled() {
	mc.called = false
}

func TestQueueCache(t *testing.T) {
	mc := &mockClient{}
	qc := NewQueueCache(mc)
	mc.setQueue(&resource.Queue{
		ID:      "1",
		Name:    "sup",
		Version: resource.NewVersion(1),
	})

	ctx := context.Background()
	q, err := qc.Get(ctx, &resource.Job{
		ID:           "1",
		QueueID:      "1",
		Version:      resource.NewVersion(1),
		QueueVersion: resource.NewVersion(1),
	})
	if err != nil {
		t.Fatal(err)
	}
	if q == nil {
		t.Fatal("expected to get a queue")
	}
	if q.ID != "1" {
		t.Fatal("expected queue id 1")
	}
	if !mc.called {
		t.Fatal("GetQueue wasn't called")
	}

	mc.resetCalled()
	q, err = qc.Get(ctx, &resource.Job{
		ID:           "1",
		QueueID:      "1",
		Version:      resource.NewVersion(1),
		QueueVersion: resource.NewVersion(1),
	})
	if err != nil {
		t.Fatal(err)
	}
	if q == nil {
		t.Fatal("expected to get a queue")
	}
	if q.ID != "1" {
		t.Fatal("expected queue id 1")
	}
	if mc.called {
		t.Fatal("GetQueue shouldn't have been called")
	}

	mc.resetCalled()
	qc.Reset()
	q, err = qc.Get(ctx, &resource.Job{
		ID:           "1",
		QueueID:      "1",
		Version:      resource.NewVersion(1),
		QueueVersion: resource.NewVersion(1),
	})
	if err != nil {
		t.Fatal(err)
	}
	if q == nil {
		t.Fatal("expected to get a queue")
	}
	if q.ID != "1" {
		t.Fatal("expected queue id 1")
	}
	if !mc.called {
		t.Fatal("GetQueue wasn't called")
	}
}
