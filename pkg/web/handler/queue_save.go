package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/internal"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
	"github.com/jeffrom/job-manager/pkg/schema"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func SaveQueue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	cfg := middleware.ConfigFromContext(ctx)
	reqLog := middleware.RequestLogFromContext(ctx)
	be := backend.FromMiddleware(ctx)
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

	claimDur := durationpb.New(0)
	if params.ClaimDuration != nil {
		claimDur = params.ClaimDuration
	}

	checkinDur := durationpb.New(0)
	if params.CheckinDuration != nil {
		checkinDur = params.CheckinDuration
	}

	now := timestamppb.New(internal.GetTimeProvider(ctx).Now())
	queue := &jobv1.Queue{
		Id:              queueID,
		Concurrency:     concurrency,
		Retries:         maxRetries,
		Labels:          params.Labels,
		Duration:        dur,
		ClaimDuration:   claimDur,
		CheckinDuration: checkinDur,
		Schema:          params.Schema,
		CreatedAt:       now,
		Unique:          params.Unique,
		V:               params.V,
	}
	savedQueue, err := be.SaveQueue(ctx, jobv1.NewQueueFromProto(queue))
	if err != nil {
		return handleBackendErrors(err, "queue", queueID)
	}
	return MarshalResponse(w, r, &apiv1.SaveQueueResponse{
		Queue: jobv1.NewQueueFromResource(savedQueue),
	})
}
