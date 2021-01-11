package handler

import (
	"net/http"

	apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
	"github.com/jeffrom/job-manager/mjob/label"
	"github.com/jeffrom/job-manager/mjob/resource"
	jobv1 "github.com/jeffrom/job-manager/mjob/resource/job/v1"
	"github.com/jeffrom/job-manager/pkg/backend"
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
	be := backend.FromMiddleware(ctx)
	queues, err := be.ListQueues(ctx, &resource.QueueListParams{
		Names:     params.Names,
		Selectors: sels,
	})
	if err != nil {
		return err
	}

	return MarshalResponse(w, r, &apiv1.ListQueuesResponse{
		Items: jobv1.NewQueuesFromResources(queues.Queues),
	})
}
