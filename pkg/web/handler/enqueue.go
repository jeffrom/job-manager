package handler

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"net/http"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/job"
	"github.com/jeffrom/job-manager/pkg/schema"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

type EnqueueJobs struct {
	// schemaCache
}

func (h *EnqueueJobs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Func(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		be := middleware.GetBackend(ctx)
		var params apiv1.EnqueueParams
		if err := UnmarshalBody(r, &params, true); err != nil {
			return err
		}

		jobs := &job.Jobs{Jobs: make([]*job.Job, len(params.Jobs))}
		ids := make([]string, len(params.Jobs))
		now := timestamppb.Now()
		for i, jobArg := range params.Jobs {
			queue, err := be.GetQueue(ctx, jobArg.Job)
			if err != nil {
				if errors.Is(err, backend.ErrNotFound) {
					return apiv1.NewNotFoundError("queue")
				}
				return err
			}

			// validate args if there is a schema
			scm, err := schema.Parse(queue)
			if err != nil {
				return err
			}

			if err := scm.ValidateArgs(ctx, jobArg.Args); err != nil {
				// fmt.Printf("VALIDERRRR: %+v\n", err)
				return err
			}

			if queue.Unique {
				unique, err := checkArgUniqueness(ctx, be, scm, jobArg.Args)
				if err != nil {
					return err
				}
				if unique {
					// return conflict error
					return apiv1.NewConflictError(queue.Id)
				}
			}

			id := job.NewID()
			jobs.Jobs[i] = &job.Job{
				Id:         id,
				Name:       jobArg.Job,
				Args:       jobArg.Args,
				Data:       jobArg.Data,
				Status:     job.StatusQueued,
				EnqueuedAt: now,
			}
			ids[i] = id
			// fmt.Printf("JOB: %+v\n", jobs.Jobs[i])
		}

		if err := be.EnqueueJobs(ctx, jobs); err != nil {
			return err
		}
		return MarshalResponse(w, r, &apiv1.EnqueueResponse{
			Jobs: ids,
		})
	})(w, r)
}

func checkArgUniqueness(ctx context.Context, be backend.Interface, scm *schema.Schema, args []*structpb.Value) (bool, error) {
	iargs := make([]interface{}, len(args))
	for i, arg := range args {
		iargs[i] = arg.AsInterface()
	}
	key, err := uniquenessKeyFromArgs(iargs)
	if err != nil {
		return false, err
	}
	found, err := be.GetSetJobKeys(ctx, []string{key})
	if err != nil {
		return false, err
	}
	if found {
		return true, nil
	}
	return false, nil
}

func uniquenessKeyFromArgs(args []interface{}) (string, error) {
	b, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return string(sum[:]), nil
}
