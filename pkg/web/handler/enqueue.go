package handler

import (
	"context"
	"crypto/sha256"
	"net/http"

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
	Func(h.handle)(w, r)
}

func (h *EnqueueJobs) handle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	be := backend.FromMiddleware(ctx)
	var params apiv1.EnqueueJobsRequest
	if err := UnmarshalBody(r, &params, true); err != nil {
		return err
	}

	resources := &resource.Jobs{Jobs: make([]*resource.Job, len(params.Jobs))}
	jobs := &jobv1.Jobs{Jobs: make([]*jobv1.Job, len(params.Jobs))}
	now := timestamppb.New(internal.GetTimeProvider(ctx).Now())
	uniqueArgQueues := make(map[string]bool)
	for i, jobArg := range params.Jobs {
		queue, err := be.GetQueue(ctx, jobArg.Job, nil)
		if err != nil {
			return handleBackendErrors(err, "queue", jobArg.Job)
		}

		if queue.Blocked {
			return handleBackendErrors(backend.ErrBlocked, "queue", jobArg.Job)
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
			uniqueArgQueues[jobArg.Job] = true
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
		resources.Jobs[i] = jobRes
		// fmt.Printf("JOB: %+v\n", jobs.Jobs[i])
	}

	res, err := be.EnqueueJobs(ctx, resources)
	if err != nil {
		return err
	}

	jobIDs := res.IDs()
	uniqueIDs, uniqueKeys, err := h.gatherUniqueArgs(res.Jobs, uniqueArgQueues)
	if err != nil {
		return err
	}
	if err := be.SetJobUniqueArgs(ctx, uniqueIDs, uniqueKeys); err != nil {
		return err
	}

	return MarshalResponse(w, r, &apiv1.EnqueueJobsResponse{
		Jobs: jobIDs,
	})
}

func (h *EnqueueJobs) gatherUniqueArgs(jobs []*resource.Job, uniqueQueues map[string]bool) ([]string, []string, error) {
	var ids []string
	var keys []string
	for _, jb := range jobs {
		if ok := uniqueQueues[jb.Name]; !ok {
			continue
		}
		ids = append(ids, jb.ID)
		key, err := jb.ArgKey()
		if err != nil {
			return nil, nil, err
		}
		keys = append(keys, key)
	}
	return ids, keys, nil
}

func checkArgUniqueness(ctx context.Context, be backend.Interface, scm *schema.Schema, args []byte) ([]string, bool, error) {
	key, err := uniquenessKeyFromArgs(args)
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

func uniquenessKeyFromArgs(b []byte) (string, error) {
	sum := sha256.Sum256(b)
	return string(sum[:]), nil
}
