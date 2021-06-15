package consumer

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/jeffrom/job-manager/mjob/logger"
	"github.com/jeffrom/job-manager/mjob/resource"
)

type worker struct {
	cfg    Config
	log    logger.Logger
	in     chan *resource.Job
	out    chan *resource.JobResult
	runner Runner
}

func newWorker(cfg Config, log logger.Logger, runner Runner, in chan *resource.Job, out chan *resource.JobResult) *worker {
	return &worker{
		cfg:    cfg,
		log:    log,
		in:     in,
		out:    out,
		runner: runner,
	}
}

func (w *worker) start(ctx context.Context) {
	// ctx, cancel := context.
	for {
		select {
		case <-ctx.Done():
			return
		case jb := <-w.in:
			if jb == nil {
				break
			}
			res, err := w.runOneJob(ctx, jb)
			w.respond(jb, res, err)
			if err != nil {
				w.log.Log(ctx, &logger.Event{Level: "error", Message: "Job failed: " + err.Error(), JobID: jb.ID, Data: res, Error: err})
			} else {
				w.log.Log(ctx, &logger.Event{Level: "info", Message: "Job complete", JobID: jb.ID, Data: res})
			}
		}
	}
}

func (w *worker) runOneJob(ctx context.Context, jb *resource.Job) (res *resource.JobResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			if rerr, ok := r.(error); ok {
				err = rerr
			} else {
				err = fmt.Errorf("consumer error: %+v", r)
			}
			w.log.Log(ctx, &logger.Event{
				Level:   "error",
				Message: fmt.Sprintf("panic: %+v\n%s", r, debug.Stack()),
			})
		}
	}()

	if jb.Duration > 0 {
		var done context.CancelFunc
		ctx, done = context.WithDeadline(context.Background(), time.Now().Add(time.Duration(jb.Duration)))
		defer done()
	}
	w.log.Log(ctx, &logger.Event{Level: "info", Message: "Starting job", JobID: jb.ID, Data: jb})
	res, err = w.runner.Run(ctx, jb)
	return res, err
}

func (w *worker) respond(jb *resource.Job, res *resource.JobResult, err error) {
	if res == nil {
		res = &resource.JobResult{}
	}
	res.JobID = jb.ID

	// TODO better error handling
	if err != nil {
		res.Error = err.Error()
		res.Status = resource.NewStatus(resource.StatusFailed)
	} else if res.Status == nil {
		res.Status = resource.NewStatus(resource.StatusComplete)
	}

	w.out <- res
}
