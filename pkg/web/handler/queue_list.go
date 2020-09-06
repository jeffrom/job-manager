package handler

import (
	"net/http"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/label"
	"github.com/jeffrom/job-manager/pkg/resource"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func ListQueues(w http.ResponseWriter, r *http.Request) error {
	var params apiv1.ListQueuesRequest
	if err := UnmarshalBody(r, &params, false); err != nil {
		return err
	}

	sels, err := label.ParseSelectorStringArray(params.Selectors)
	if err != nil {
		return err
	}

	ctx := r.Context()
	be := middleware.GetBackend(ctx)
	queues, err := be.ListQueues(ctx, &resource.QueueListParams{
		Names:     params.Names,
		Selectors: sels,
	})
	if err != nil {
		return err
	}

	return MarshalResponse(w, r, &apiv1.ListQueuesResponse{
		Data: &jobv1.Queues{
			Queues: jobv1.NewQueuesFromResources(queues.Queues),
		},
	})
}
