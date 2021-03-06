package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
	jobv1 "github.com/jeffrom/job-manager/mjob/resource/job/v1"
	"github.com/jeffrom/job-manager/mjob/schema"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/internal"
	"github.com/jeffrom/job-manager/pkg/logger"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func SaveQueue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	cfg := middleware.ConfigFromContext(ctx)
	reqLog := logger.RequestLogFromContext(ctx)
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

	maxRetries := int32(cfg.DefaultMaxRetries)
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

	var boFactor float32 = 2.0
	if params.BackoffFactor > 0 {
		boFactor = params.BackoffFactor
	}

	boInitial := durationpb.New(1 * time.Second)
	if params.BackoffInitialDuration != nil {
		boInitial = params.BackoffInitialDuration
	}

	boMax := durationpb.New(20 * time.Minute)
	if params.BackoffMaxDuration != nil {
		boMax = params.BackoffMaxDuration
	}

	now := timestamppb.New(internal.GetTimeProvider(ctx).Now())
	queue := &jobv1.Queue{
		Id:                     queueID,
		Retries:                maxRetries,
		Labels:                 params.Labels,
		Duration:               dur,
		ClaimDuration:          claimDur,
		CheckinDuration:        checkinDur,
		Schema:                 params.Schema,
		CreatedAt:              now,
		Unique:                 params.Unique,
		V:                      params.V,
		BackoffFactor:          boFactor,
		BackoffInitialDuration: boInitial,
		BackoffMaxDuration:     boMax,
	}
	savedQueue, err := be.SaveQueue(ctx, jobv1.NewQueueFromProto(queue))
	if err != nil {
		return handleBackendErrors(err, "queue", queueID)
	}
	return MarshalResponse(w, r, &apiv1.SaveQueueResponse{
		Queue: jobv1.NewQueueFromResource(savedQueue),
	})
}
