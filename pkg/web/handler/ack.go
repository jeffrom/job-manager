package handler

import (
	"net/http"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/job"
	"github.com/jeffrom/job-manager/pkg/schema"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func Ack(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := middleware.GetBackend(ctx)
	var params apiv1.AckParams
	if err := UnmarshalBody(r, &params, true); err != nil {
		return err
	}

	results := &job.Acks{Acks: make([]*job.Ack, len(params.Acks))}
	for i, ackParam := range params.Acks {
		id := ackParam.Id
		jobData, err := be.GetJobByID(ctx, id)
		if err != nil {
			return err
		}
		queue, err := be.GetQueue(ctx, jobData.Name)
		if err != nil {
			return err
		}
		scm, err := schema.Parse(queue)
		if err != nil {
			return err
		}
		if err := scm.ValidateResult(ctx, ackParam.Data); err != nil {
			return err
		}

		results.Acks[i] = &job.Ack{
			Id:     id,
			Status: ackParam.Status,
			Data:   ackParam.Data,
		}
	}

	if err := be.AckJobs(ctx, results); err != nil {
		return err
	}
	return nil
}
