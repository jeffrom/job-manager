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

func DequeueJobs(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := backend.FromMiddleware(ctx)
	queueID := chi.URLParam(r, "queueID")
	var params apiv1.DequeueRequest
	if err := UnmarshalBody(r, &params, queueID == ""); err != nil {
		return err
	}
	// TODO error if queueID url param is set and there are more than one queue
	if queueID != "" && len(params.Queues) == 0 {
		params.Queues = []string{queueID}
	}
	var num int = 1
	if params.Num > 0 {
		num = int(params.Num)
	}

	claims, err := label.ParseClaims(params.Claims)
	if err != nil {
		return err
	}

	for _, qName := range params.Queues {
		_, err = be.GetQueue(ctx, qName)
		if err != nil {
			return err
		}
	}

	listOpts := &resource.JobListParams{Claims: claims}
	jobs, err := be.DequeueJobs(ctx, num, listOpts)
	if err != nil {
		return err
	}

	jobsResp, err := jobv1.NewJobsFromResources(jobs.Jobs)
	if err != nil {
		return err
	}
	return MarshalResponse(w, r, &apiv1.DequeueResponse{
		Jobs: &jobv1.Jobs{
			Jobs: jobsResp,
		},
	})
}
