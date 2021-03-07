package handler

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"net/http"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
	"github.com/jeffrom/job-manager/mjob/label"
	"github.com/jeffrom/job-manager/mjob/resource"
	jobv1 "github.com/jeffrom/job-manager/mjob/resource/job/v1"
	"github.com/jeffrom/job-manager/mjob/schema"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/internal"
)

type EnqueueJobs struct {
	// schemaCache
}

func (h *EnqueueJobs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Func(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		be := backend.FromMiddleware(ctx)
		var params apiv1.EnqueueJobsRequest
		if err := UnmarshalBody(r, &params, true); err != nil {
			return err
		}

		resources := &resource.Jobs{Jobs: make([]*resource.Job, len(params.Jobs))}
		jobs := &jobv1.Jobs{Jobs: make([]*jobv1.Job, len(params.Jobs))}
		now := timestamppb.New(internal.GetTimeProvider(ctx).Now())
		for i, jobArg := range params.Jobs {
			queue, err := be.GetQueue(ctx, jobArg.Job, nil)
			if err != nil {
				return handleBackendErrors(err, "queue", jobArg.Job)
			}

			// validate args if there is a schema
			scm, err := schema.Parse(queue.SchemaRaw)
			if err != nil {
				return err
			}

			if err := scm.ValidateArgs(ctx, jobArg.Args); err != nil {
				return handleSchemaErrors(err, "job", "", "invalid job arguments")
			}

			if queue.Unique {
				ids, unique, err := checkArgUniqueness(ctx, be, scm, jobArg.Args)
				if err != nil {
					return err
				}
				if unique {
					// return conflict error
					return resource.NewUnprocessableEntityError("queue", queue.Name, "A job with matching arguments is executing", ids)
				}
			}

			var claims label.Claims
			if jobArg.Data != nil {
				claims, err = label.ParseClaims(jobArg.Data.Claims)
				if err != nil {
					return err
				}
			}

			jb := &jobv1.Job{
				Name:       jobArg.Job,
				Args:       jobArg.Args,
				Data:       jobArg.Data,
				Status:     jobv1.StatusQueued,
				EnqueuedAt: now,
			}
			jobs.Jobs[i] = jb

			jobRes := jobv1.NewJobFromProto(jb, claims)
			if err := jobRes.Populate(); err != nil {
				return err
			}
			resources.Jobs[i] = jobRes
			// fmt.Printf("JOB: %+v\n", jobs.Jobs[i])
		}

		res, err := be.EnqueueJobs(ctx, resources)
		if err != nil {
			return err
		}

		jobIDs := res.IDs()
		jobKeys, err := res.ArgKeys()
		if err != nil {
			return err
		}
		if err := be.SetJobUniqueArgs(ctx, jobIDs, jobKeys); err != nil {
			return err
		}
		return MarshalResponse(w, r, &apiv1.EnqueueJobsResponse{
			Jobs: jobIDs,
		})
	})(w, r)
}

func checkArgUniqueness(ctx context.Context, be backend.Interface, scm *schema.Schema, args []*structpb.Value) ([]string, bool, error) {
	iargs := make([]interface{}, len(args))
	for i, arg := range args {
		iargs[i] = arg.AsInterface()
	}
	key, err := uniquenessKeyFromArgs(iargs)
	if err != nil {
		return nil, false, err
	}
	ids, found, err := be.GetJobUniqueArgs(ctx, []string{key})
	if err != nil {
		return nil, false, err
	}
	if found {
		return ids, true, nil
	}
	return nil, false, nil
}

func uniquenessKeyFromArgs(args []interface{}) (string, error) {
	b, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return string(sum[:]), nil
}
