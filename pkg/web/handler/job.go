package handler

import (
	"net/http"

	"github.com/go-chi/chi"

	apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
	"github.com/jeffrom/job-manager/mjob/label"
	"github.com/jeffrom/job-manager/mjob/resource"
	jobv1 "github.com/jeffrom/job-manager/mjob/resource/job/v1"
	"github.com/jeffrom/job-manager/pkg/backend"
)

func ListJobs(w http.ResponseWriter, r *http.Request) error {
	var params apiv1.ListJobsRequest
	if err := UnmarshalBody(r, &params, false); err != nil {
		return err
	}

	sels, err := label.ParseSelectorStringArray(params.Selector)
	if err != nil {
		return err
	}

	claims, err := label.ParseClaims(params.Claim)
	if err != nil {
		return err
	}

	ctx := r.Context()
	be := backend.FromMiddleware(ctx)

	// fmt.Println("status", params.Statuses)
	resourceParams := &resource.JobListParams{
		Names:       params.Name,
		Statuses:    resource.StatusesFromStrings(params.Status...),
		Selectors:   sels,
		Claims:      claims,
		NoUnclaimed: params.NoUnclaimed,
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
