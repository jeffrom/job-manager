package integration_test

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"testing"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/consumer"
	"github.com/jeffrom/job-manager/mjob/logger"
	"github.com/jeffrom/job-manager/mjob/resource"
	bememory "github.com/jeffrom/job-manager/pkg/backend/mem"
	"github.com/jeffrom/job-manager/pkg/testenv"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func BenchmarkMemory(b *testing.B) {
	b.SetParallelism(1)
	n := 1
	if env := os.Getenv("N"); env != "" {
		envN, err := strconv.ParseInt(env, 10, 64)
		if err != nil {
			b.Fatal(err)
		}
		n = int(envN)
	}
	cpus := runtime.NumCPU()
	b.Logf("consumer concurrency ($N): %d * %d cpus", n, cpus)
	cfg := middleware.NewConfig()
	be := bememory.New()
	srv := testenv.NewTestControllerServer(b, cfg, be)
	srv.Start()
	defer srv.Close()

	ctx := context.Background()
	c := testenv.NewTestClient(b, srv)
	_, err := c.SaveQueue(ctx, "benchmem", client.SaveQueueOpts{MaxRetries: 0})
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		if _, err := c.EnqueueJob(ctx, "benchmem"); err != nil {
			b.Fatal(err)
		}
	}

	br := newBenchRunner()
	cons := consumer.New(c, br,
		consumer.WithConfig(consumer.Config{Concurrency: cpus * n}),
		consumer.WithLogger(logger.Error),
	)
	defer cons.Stop()
	consErrC := make(chan error)
	go func() {
		if err := cons.Run(ctx); err != nil {
			consErrC <- err
		}
		consErrC <- nil
	}()

	b.ResetTimer()
	// enqueue a job and wait for it to complete, using a channel for
	// communication between the consumer and the benchmark goroutine
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
	b.StopTimer()

	cons.Stop()
	select {
	case err := <-consErrC:
		if err != nil {
			b.Fatal(err)
		}
	default:
	}
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
	if len(job.ArgsRaw) == 0 {
		return nil, nil
	}
	args := []string{}
	if err := json.Unmarshal(job.ArgsRaw, &args); err != nil {
		return nil, err
	}
	ch := br.pop(args[0])
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
