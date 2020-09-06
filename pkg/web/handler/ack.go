package handler

import (
	"context"
	"net/http"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/backend"
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
		scm, err := job.ParseSchema(queue)
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
	if err := deleteArgUniqueness(ctx, be, results.Acks); err != nil {
		return err
	}
	return nil
}

func deleteArgUniqueness(ctx context.Context, be backend.Interface, acks []*job.Ack) error {
	var keys []string
	for _, ack := range acks {
		if !job.IsComplete(ack.Status) {
			continue
		}

		job, err := be.GetJobByID(ctx, ack.Id)
		if err != nil {
			return err
		}
		iargs := make([]interface{}, len(job.Args))
		for i, arg := range job.Args {
			iargs[i] = arg.AsInterface()
		}
		ukey, err := uniquenessKeyFromArgs(iargs)
		if err != nil {
			return err
		}
		keys = append(keys, ukey)
	}
	return be.DeleteJobKeys(ctx, keys)
}
