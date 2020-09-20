package handler

import (
	"net/http"

	"github.com/go-chi/chi"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/label"
	"github.com/jeffrom/job-manager/pkg/resource"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
)

func ListJobs(w http.ResponseWriter, r *http.Request) error {
	var params jobv1.JobListParams
	if err := UnmarshalBody(r, &params, false); err != nil {
		return err
	}

	sels, err := label.ParseSelectorStringArray(params.Selectors)
	if err != nil {
		return err
	}

	ctx := r.Context()
	be := backend.FromMiddleware(ctx)

	resourceParams := &resource.JobListParams{
		Names:     params.Names,
		Selectors: sels,
	}
	jobs, err := be.ListJobs(ctx, 20, resourceParams)
	if err != nil {
		return err
	}

	respJobs, err := jobv1.NewJobsFromResources(jobs.Jobs)
	if err != nil {
		return err
	}
	return MarshalResponse(w, r, &apiv1.ListJobsResponse{
		Items: respJobs,
	})
}

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
