package handler

import (
	"net/http"

	apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
	"github.com/jeffrom/job-manager/pkg/backend"
)

func Stats(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := backend.FromMiddleware(ctx)
	stats, err := be.Stats(ctx)
	if err != nil {
		return err
	}

	return MarshalResponse(w, r, &apiv1.StatsResponse{
		Queued:               stats.Queued,
		Running:              stats.Running,
		Complete:             stats.Complete,
		Dead:                 stats.Dead,
		Cancelled:            stats.Cancelled,
		Invalid:              stats.Invalid,
		Failed:               stats.Failed,
		LongestUnstartedSecs: stats.LongestUnstarted,
	})
}
