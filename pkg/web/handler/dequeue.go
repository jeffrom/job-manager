package handler

import (
	"net/http"

	"github.com/go-chi/chi"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/job"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func DequeueJobs(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := middleware.GetBackend(ctx)
	queueID := chi.URLParam(r, "queueID")
	var params apiv1.DequeueParams
	if err := UnmarshalBody(r, &params, queueID == ""); err != nil {
		return err
	}
	if queueID != "" {
		params.Job = queueID
	}

	_, err := be.GetQueue(ctx, queueID)
	if err != nil {
		return err
	}

	var num int = 1
	if params.Num > 0 {
		num = int(params.Num)
	}
	listOpts := &job.JobListParams{Statuses: []job.Status{job.StatusQueued}}
	jobs, err := be.DequeueJobs(ctx, num, listOpts)
	if err != nil {
		return err
	}
	return MarshalResponse(w, r, &apiv1.DequeueResponse{Jobs: jobs})
}
