package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"google.golang.org/protobuf/proto"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/job"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func DequeueJobs(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := middleware.GetBackend(ctx)
	queueName := chi.URLParam(r, "queueName")
	var params apiv1.DequeueParams
	if err := UnmarshalBody(r, &params, queueName == ""); err != nil {
		return err
	}
	if queueName != "" {
		params.Job = proto.String(queueName)
	}

	_, err := be.GetQueue(ctx, queueName)
	if err != nil {
		return err
	}

	var num int = 1
	if params.Num != nil {
		num = int(*params.Num)
	}
	listOpts := &job.ListOpts{Statuses: []job.Status{job.StatusQueued}}
	jobs, err := be.DequeueJobs(ctx, num, listOpts)
	if err != nil {
		return err
	}
	return MarshalResponse(w, r, jobs)
}
