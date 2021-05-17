package integration

import (
	"context"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/consumer"
	"github.com/jeffrom/job-manager/mjob/logger"
	"github.com/jeffrom/job-manager/mjob/resource"
	"github.com/jeffrom/job-manager/pkg/backend/mem"
	srvlogger "github.com/jeffrom/job-manager/pkg/logger"
	"github.com/jeffrom/job-manager/pkg/testenv"
	"github.com/jeffrom/job-manager/pkg/web"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
	"github.com/rs/zerolog"
)

func TestMemoryCounter(t *testing.T) {
	t.Skip("just for funsies")
	if testing.Short() {
		t.Skip("-short")
	}
	n := 10
	if env := os.Getenv("N"); env != "" {
		envN, err := strconv.ParseInt(env, 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		n = int(envN)
	}
	cpus := runtime.NumCPU()
	t.Logf("consumer concurrency ($N): %d (%d cpus)", n, cpus)

	cfg := middleware.NewConfig()
	cfg.Logger = &srvlogger.Logger{Disabled: true, Logger: zerolog.Nop()}
	cfg.ResetLogOutput(ioutil.Discard)
	be := mem.New()
	h, err := web.NewControllerRouter(cfg, be)
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewUnstartedServer(h)
	t.Logf("Started job-controller server with backend %T at address: %s", be, srv.Listener.Addr())
	srv.Start()
	defer srv.Close()

	ctx := context.Background()
	c := testenv.NewTestClient(t, srv)
	_, err = c.SaveQueue(ctx, "memcounter", client.SaveQueueOpts{MaxRetries: 0})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 100000; i++ {
		if _, err := c.EnqueueJob(ctx, "memcounter"); err != nil {
			t.Fatal(err)
		}
	}

	cr := &counterRunner{done: make(chan error)}
	cons := consumer.New(c, cr,
		consumer.WithConfig(consumer.Config{Concurrency: n}),
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
	startedAt := time.Now()

	// for i := 0; i < 99000; i++ {
	// 	if _, err := c.EnqueueJob(ctx, "memcounter"); err != nil {
	// 		t.Fatal(err)
	// 	}
	// }

	if err := <-cr.done; err != nil {
		t.Fatal(err)
	}
	dur := time.Since(startedAt)

	cons.Stop()
	count := atomic.LoadInt64(&cr.n)
	t.Logf("counter: %d", count)
	t.Logf("took: %s (%s/job, %2f jobs/s)", dur, dur/time.Duration(count), float64(count)/dur.Seconds())
	select {
	case err := <-consErrC:
		if err != nil {
			t.Fatal(err)
		}
	default:
	}
}

type counterRunner struct {
	n    int64
	done chan error
}

func (cr *counterRunner) Run(ctx context.Context, job *resource.Job) (*resource.JobResult, error) {
	n := atomic.AddInt64(&cr.n, 1)
	if n >= 100000 {
		cr.done <- nil
	}
	// fmt.Println("SUP", n)
	return nil, nil
}
