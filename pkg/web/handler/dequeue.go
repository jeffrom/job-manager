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

func DequeueJobs(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := backend.FromMiddleware(ctx)
	queueID := chi.URLParam(r, "queueID")
	var params apiv1.DequeueJobsRequest
	if err := UnmarshalBody(r, &params, false); err != nil {
		return err
	}
	// TODO error if queueID url param is set and there are more than one queue
	if queueID != "" && len(params.Queues) == 0 {
		params.Queues = []string{queueID}
	}
	num := 1
	if params.Num > 0 {
		num = int(params.Num)
	}

	claims, err := label.ParseClaims(params.Claims)
	if err != nil {
		return err
	}

	for _, qName := range params.Queues {
		_, err = be.GetQueue(ctx, qName, nil)
		if err != nil {
			return err
		}
	}

	// fmt.Println("params:", params.Claims, "parsed:", claims)
	listOpts := &resource.JobListParams{
		Queues: params.Queues,
		Claims: claims,
	}
	jobs, err := be.DequeueJobs(ctx, num, listOpts)
	if err != nil {
		return err
	}

	jobsResp, err := jobv1.NewJobsFromResources(jobs.Jobs)
	if err != nil {
		return err
	}
	return MarshalResponse(w, r, &apiv1.DequeueJobsResponse{
		Items: jobsResp,
	})
}
