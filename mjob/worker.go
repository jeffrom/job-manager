package mjob

import (
	"context"

	"github.com/jeffrom/job-manager/pkg/resource"
)

type consumerWorker struct {
	cfg    ConsumerConfig
	in     chan *resource.Job
	out    chan *resource.JobResult
	runner Runner
}

func newWorker(cfg ConsumerConfig, in chan *resource.Job, out chan *resource.JobResult) *consumerWorker {
	return &consumerWorker{
		cfg: cfg,
		in:  in,
		out: out,
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
				return
			}
			w.respond(w.runner.Run(ctx, jb))
		}
	}
}

func (w *consumerWorker) respond(res *resource.JobResult, err error) {
	if res == nil {
		res = &resource.JobResult{}
	}
	if err != nil {
		res.Error = err.Error()
	}

	w.out <- res
}
