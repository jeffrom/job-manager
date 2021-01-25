package mjob

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/jeffrom/job-manager/mjob/resource"
)

type consumerWorker struct {
	cfg    ConsumerConfig
	logger Logger
	in     chan *resource.Job
	out    chan *resource.JobResult
	runner Runner
}

func newWorker(cfg ConsumerConfig, logger Logger, runner Runner, in chan *resource.Job, out chan *resource.JobResult) *consumerWorker {
	return &consumerWorker{
		cfg:    cfg,
		logger: logger,
		in:     in,
		out:    out,
		runner: runner,
	}
}

func (w *consumerWorker) start(ctx context.Context) {
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
				w.logger.Log(ctx, &LogEvent{Level: "error", Message: "Job failed: " + err.Error(), JobID: jb.ID, Data: res, Error: err})
			} else {
				w.logger.Log(ctx, &LogEvent{Level: "info", Message: "Job complete", JobID: jb.ID, Data: res})
			}
		}
	}
}

func (w *consumerWorker) runOneJob(ctx context.Context, jb *resource.Job) (res *resource.JobResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			if rerr, ok := r.(error); ok {
				err = rerr
			} else {
				err = errors.New(fmt.Sprintf("consumer error: %+v", r))
			}
			w.logger.Log(ctx, &LogEvent{
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
	w.logger.Log(ctx, &LogEvent{Level: "info", Message: "Starting job", JobID: jb.ID, Data: jb})
	res, err = w.runner.Run(ctx, jb)
	return res, err
}

func (w *consumerWorker) respond(jb *resource.Job, res *resource.JobResult, err error) {
	if res == nil {
		res = &resource.JobResult{}
	}
	res.JobID = jb.ID

	if err != nil {
		res.Error = err.Error()
		// TODO better error handling
		res.Status = resource.NewStatus(resource.StatusFailed)
	} else if res.Status == nil {
		res.Status = resource.NewStatus(resource.StatusComplete)
	}

	w.out <- res
}
