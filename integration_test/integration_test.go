package integration_test

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/jeffrom/job-manager/mjob"
	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/resource"
	"github.com/jeffrom/job-manager/pkg/backend/bememory"
	"github.com/jeffrom/job-manager/pkg/testenv"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func BenchmarkMemory(b *testing.B) {
	cfg := middleware.NewConfig()
	be := bememory.New()
	srv := testenv.NewTestControllerServer(b, cfg, be)
	srv.Start()
	defer srv.Close()

	ctx := context.Background()
	c := testenv.NewTestClient(b, srv)
	_, err := c.SaveQueue(ctx, "benchmem", client.SaveQueueOpts{})
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		if _, err := c.EnqueueJob(ctx, "benchmem"); err != nil {
			b.Fatal(err)
		}
	}

	br := newBenchRunner()
	cons := mjob.NewConsumer(c, br, mjob.ConsumerWithLogger(&mjob.NilLogger{}))
	defer cons.Stop()
	go func() {
		if err := cons.Run(ctx); err != nil {
			b.Fatal(err)
		}
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			chanID := strconv.FormatInt(rand.Int63(), 10)
			// fmt.Println("chanID", chanID)
			ch := make(chan struct{})
			br.add(chanID, ch)
			_, err := c.EnqueueJob(ctx, "benchmem", chanID)
			<-ch
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

type benchRunner struct {
	mu sync.Mutex
	m  map[string]chan struct{}
}

func (br *benchRunner) add(k string, ch chan struct{}) {
	br.mu.Lock()
	defer br.mu.Unlock()

	if _, ok := br.m[k]; ok {
		panic("double write to sync map: " + k)
	}

	br.m[k] = ch
}

func (br *benchRunner) pop(k string) chan struct{} {
	br.mu.Lock()
	defer br.mu.Unlock()

	res := br.m[k]
	delete(br.m, k)
	return res
}

func newBenchRunner() *benchRunner {
	return &benchRunner{
		m: make(map[string]chan struct{}, 1),
	}
}

func (br *benchRunner) Run(ctx context.Context, job *resource.Job) (*resource.JobResult, error) {
	if len(job.Args) == 0 {
		return nil, nil
	}
	ch := br.pop(job.Args[0].(string))
	// fmt.Println("job arg 0:", job.Args[0], ch, job.Attempt)
	if ch == nil {
		return nil, errors.New("no channel")
	}
	defer func() {
		ch <- struct{}{}
	}()
	return nil, nil
}

// func TestIntegrationServerStart(t *testing.T) {
// 	srv := testenv.NewTestControllerServer(t, nil)
// 	srv.Start()
// 	defer srv.Close()
// }

// func TestIntegrationClientConnect(t *testing.T) {
// 	srv := testenv.NewTestControllerServer(t, nil)
// 	c := testenv.NewTestClient(t, srv)
// 	srv.Start()
// 	defer srv.Close()

// 	if err := c.Ping(context.Background()); err != nil {
// 		t.Fatal(err)
// 	}
// }
