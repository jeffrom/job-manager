package mjob

import (
	"context"

	"github.com/jeffrom/job-manager/mjob/resource"
)

type consumerWorker struct {
	cfg    ConsumerConfig
	in     chan *resource.Job
	out    chan *resource.JobResult
	runner Runner
}

func newWorker(cfg ConsumerConfig, runner Runner, in chan *resource.Job, out chan *resource.JobResult) *consumerWorker {
	return &consumerWorker{
		cfg:    cfg,
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
			res, err := w.runner.Run(ctx, jb)
			w.respond(jb, res, err)
		}
	}
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
