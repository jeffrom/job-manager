package handler

import (
	"net/http"

	"github.com/go-chi/chi"

	apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
	jobv1 "github.com/jeffrom/job-manager/mjob/resource/job/v1"
	"github.com/jeffrom/job-manager/pkg/backend"
)

func DeleteQueue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := backend.FromMiddleware(ctx)
	queueID := chi.URLParam(r, "queueID")

	if err := be.DeleteQueues(ctx, []string{queueID}); err != nil {
		return handleBackendErrors(err, "queue", queueID)
	}
	return MarshalResponse(w, r, &apiv1.DeleteQueueResponse{Ok: true})
}

func PauseQueue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := backend.FromMiddleware(ctx)
	queueID := chi.URLParam(r, "queueID")

	if err := be.PauseQueues(ctx, []string{queueID}); err != nil {
		return handleBackendErrors(err, "queue", queueID)
	}
	return MarshalResponse(w, r, &apiv1.PauseQueueResponse{Ok: true})
}

func UnpauseQueue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := backend.FromMiddleware(ctx)
	queueID := chi.URLParam(r, "queueID")

	if err := be.UnpauseQueues(ctx, []string{queueID}); err != nil {
		return handleBackendErrors(err, "queue", queueID)
	}
	return MarshalResponse(w, r, &apiv1.PauseQueueResponse{Ok: true})
}

func BlockQueue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := backend.FromMiddleware(ctx)
	queueID := chi.URLParam(r, "queueID")

	if err := be.BlockQueues(ctx, []string{queueID}); err != nil {
		return handleBackendErrors(err, "queue", queueID)
	}
	return MarshalResponse(w, r, &apiv1.BlockQueueResponse{Ok: true})
}

func UnblockQueue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := backend.FromMiddleware(ctx)
	queueID := chi.URLParam(r, "queueID")

	if err := be.UnblockQueues(ctx, []string{queueID}); err != nil {
		return handleBackendErrors(err, "queue", queueID)
	}
	return MarshalResponse(w, r, &apiv1.UnblockQueueResponse{Ok: true})
}

func GetQueueByID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := backend.FromMiddleware(ctx)
	queueID := chi.URLParam(r, "queueID")

	queue, err := be.GetQueue(ctx, queueID, nil)
	if err != nil {
		return handleBackendErrors(err, "queue", queueID)
	}
	// fmt.Printf("kewl %+v\n", queue)
	return MarshalResponse(w, r, &apiv1.GetQueueResponse{
		Data: jobv1.NewQueueFromResource(queue),
	})
}

func GetQueueByJobID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := backend.FromMiddleware(ctx)
	jobID := chi.URLParam(r, "jobID")
	jobData, err := be.GetJobByID(ctx, jobID, nil)
	if err != nil {
		return handleBackendErrors(err, "job", jobID)
	}

	queue, err := be.GetQueue(ctx, jobData.Name, nil)
	if err != nil {
		return handleBackendErrors(err, "queue", jobData.Name)
	}
	return MarshalResponse(w, r, &apiv1.GetQueueResponse{
		Data: jobv1.NewQueueFromResource(queue),
	})
}
