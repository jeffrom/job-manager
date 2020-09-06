package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
	"github.com/jeffrom/job-manager/pkg/schema"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func SaveQueue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	cfg := middleware.ConfigFromContext(ctx)
	reqLog := middleware.RequestLogFromContext(ctx)
	be := middleware.GetBackend(ctx)
	queueID := chi.URLParam(r, "queueID")
	reqLog.Str("queue", queueID)

	var params apiv1.SaveQueueParamArgs
	if err := UnmarshalBody(r, &params, false); err != nil {
		return err
	}

	if err := schema.ValidateSchema(ctx, params.Schema); err != nil {
		return err
	}

	var concurrency int32 = int32(cfg.DefaultConcurrency)
	if conc := params.Concurrency; conc > 0 {
		concurrency = conc
	}
	var maxRetries int32 = int32(cfg.DefaultMaxRetries)
	if mr := params.MaxRetries; mr > 0 {
		maxRetries = mr
	}

	dur := durationpb.New(cfg.DefaultMaxJobTimeout)
	if d := params.Duration; d != nil {
		dur = d
	}

	now := timestamppb.Now()
	queue := &jobv1.Queue{
		Id:          queueID,
		Concurrency: concurrency,
		Retries:     maxRetries,
		Labels:      params.Labels,
		Duration:    dur,
		Schema:      params.Schema,
		CreatedAt:   now,
		Unique:      params.Unique,
		V:           params.V,
	}
	if err := be.SaveQueue(ctx, queue); err != nil {
		return handleBackendErrors(err, "queue", queueID)
	}
	return MarshalResponse(w, r, &apiv1.SaveQueueResponse{Queue: queue})
}

func DeleteQueue(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func ListQueues(w http.ResponseWriter, r *http.Request) error {
	var params apiv1.ListQueuesRequest
	if err := UnmarshalBody(r, &params, false); err != nil {
		return err
	}

	ctx := r.Context()
	be := middleware.GetBackend(ctx)
	queues, err := be.ListQueues(ctx, &jobv1.QueueListParams{
		Names:     params.Names,
		Selectors: params.Selectors,
	})
	if err != nil {
		return err
	}

	return MarshalResponse(w, r, &apiv1.ListQueuesResponse{Data: queues})
}

func GetQueueByID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := middleware.GetBackend(ctx)
	queueID := chi.URLParam(r, "queueID")

	queue, err := be.GetQueue(ctx, queueID)
	if err != nil {
		return handleBackendErrors(err, "queue", queueID)
	}
	return MarshalResponse(w, r, &apiv1.GetQueueResponse{Data: queue})
}

func GetQueueByJobID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := middleware.GetBackend(ctx)
	jobID := chi.URLParam(r, "jobID")
	jobData, err := be.GetJobByID(ctx, jobID)
	if err != nil {
		return handleBackendErrors(err, "job", jobID)
	}

	queue, err := be.GetQueue(ctx, jobData.Name)
	if err != nil {
		return handleBackendErrors(err, "queue", jobData.Name)
	}
	return MarshalResponse(w, r, &apiv1.GetQueueResponse{Data: queue})
}
