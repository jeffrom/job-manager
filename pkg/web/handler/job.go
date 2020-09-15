package handler

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/jeffrom/job-manager/pkg/backend"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
)

func GetJobByID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := backend.FromMiddleware(ctx)
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
