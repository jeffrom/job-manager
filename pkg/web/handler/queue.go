package handler

import (
	"net/http"

	"github.com/go-chi/chi"

	apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
	jobv1 "github.com/jeffrom/job-manager/mjob/resource/job/v1"
	"github.com/jeffrom/job-manager/pkg/backend"
)

func DeleteQueue(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func GetQueueByID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := backend.FromMiddleware(ctx)
	queueID := chi.URLParam(r, "queueID")

	queue, err := be.GetQueue(ctx, queueID)
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
	jobData, err := be.GetJobByID(ctx, jobID)
	if err != nil {
		return handleBackendErrors(err, "job", jobID)
	}

	queue, err := be.GetQueue(ctx, jobData.Name)
	if err != nil {
		return handleBackendErrors(err, "queue", jobData.Name)
	}
	return MarshalResponse(w, r, &apiv1.GetQueueResponse{
		Data: jobv1.NewQueueFromResource(queue),
	})
}
