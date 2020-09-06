package handler

import (
	"net/http"

	"github.com/go-chi/chi"

	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func GetJobByID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := middleware.GetBackend(ctx)
	jobID := chi.URLParam(r, "jobID")

	job, err := be.GetJobByID(ctx, jobID)
	if err != nil {
		return err
	}
	respJob, err := jobv1.NewJobFromResource(job)
	if err != nil {
		return err
	}
	return MarshalResponse(w, r, respJob)
}
