package handler

import (
	"context"
	"net/http"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/resource"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
	"github.com/jeffrom/job-manager/pkg/schema"
)

func Ack(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := backend.FromMiddleware(ctx)
	var params apiv1.AckJobsRequest
	if err := UnmarshalBody(r, &params, true); err != nil {
		return err
	}

	resources := &resource.Acks{Acks: make([]*resource.Ack, len(params.Acks))}
	results := &jobv1.Acks{Acks: make([]*jobv1.Ack, len(params.Acks))}
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
		scm, err := schema.Parse(queue.SchemaRaw)
		if err != nil {
			return err
		}
		if err := scm.ValidateResult(ctx, ackParam.Data); err != nil {
			return err
		}

		ack := &jobv1.Ack{
			Id:     id,
			Status: ackParam.Status,
			Data:   ackParam.Data,
			Error:  ackParam.Error,
		}
		results.Acks[i] = ack
		ackRes := jobv1.AckFromProto(ack)
		resources.Acks[i] = ackRes
	}

	if err := be.AckJobs(ctx, resources); err != nil {
		return handleBackendErrors(err, "ack", "")
	}
	if err := deleteArgUniqueness(ctx, be, resources.Acks); err != nil {
		return err
	}
	return nil
}

func deleteArgUniqueness(ctx context.Context, be backend.Interface, acks []*resource.Ack) error {
	var keys []string
	for _, ack := range acks {
		if !resource.StatusIsAttempted(ack.Status) {
			continue
		}

		jobData, err := be.GetJobByID(ctx, ack.JobID)
		if err != nil {
			return err
		}
		iargs := make([]interface{}, len(jobData.Args))
		for i, arg := range jobData.Args {
			iargs[i] = arg
		}
		ukey, err := uniquenessKeyFromArgs(iargs)
		if err != nil {
			return err
		}
		keys = append(keys, ukey)
	}
	return be.DeleteJobKeys(ctx, keys)
}
