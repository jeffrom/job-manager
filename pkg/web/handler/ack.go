package handler

import (
	"net/http"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/job"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func Ack(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := middleware.GetBackend(ctx)
	var params apiv1.AckParams
	if err := UnmarshalBody(r, &params, true); err != nil {
		return err
	}

	results := &job.Results{Results: make([]*job.Result, len(params.Acks))}
	for i, ackParam := range params.Acks {
		results.Results[i] = ackParam.Result
	}

	if err := be.AckJobs(ctx, results); err != nil {
		return err
	}
	return nil
}
