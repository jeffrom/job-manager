package handler

import (
	"errors"
	"net/http"
	"time"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/job"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func EnqueueJobs(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := middleware.GetBackend(ctx)
	var params apiv1.EnqueueParams
	if err := UnmarshalBody(r, &params, true); err != nil {
		return err
	}

	jobs := &job.Jobs{
		Jobs: make([]*job.Job, len(params.Jobs)),
	}
	now := timestamppb.Now()
	for i, jobArg := range params.Jobs {
		queue, err := be.GetQueue(ctx, jobArg.Job)
		if err != nil {
			if errors.Is(err, backend.ErrNotFound) {
				return apiv1.NewNotFoundError("queue")
			}
			return err
		}

		var maxRetries int32
		if jobArg.Retries != nil {
			maxRetries = *jobArg.Retries
		} else if queue != nil && queue.MaxRetries > 0 {
			maxRetries = queue.MaxRetries
		}

		var dur *durationpb.Duration
		if jobArg.Duration != nil {
			dur = jobArg.Duration
		} else if d := queue.Duration; d != nil && (d.Seconds > 0 || d.Nanos > 0) {
			dur = queue.Duration
		} else {
			dur = durationpb.New(10 * time.Minute)
		}

		id := job.NewID()
		jobs.Jobs[i] = &job.Job{
			Id:         id,
			Name:       jobArg.Job,
			Args:       jobArg.Args,
			Data:       jobArg.Data,
			Duration:   dur,
			MaxRetries: maxRetries,
			Status:     job.StatusQueued,
			EnqueuedAt: now,
		}
	}

	if err := be.EnqueueJobs(ctx, jobs); err != nil {
		return err
	}
	return MarshalResponse(w, r, jobs)
}
