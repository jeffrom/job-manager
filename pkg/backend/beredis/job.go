package beredis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/jeffrom/job-manager/pkg/internal"
	"github.com/jeffrom/job-manager/pkg/resource"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
)

const streamKey = "mjob:jobs"

func (be *RedisBackend) EnqueueJobs(ctx context.Context, jobs *resource.Jobs) (*resource.Jobs, error) {
	// first check everything has a queue
	res, err := be.queuesForJobs(ctx, jobs.Jobs)
	if err != nil {
		return nil, err
	}
	queues := res.ToMap()

	now := internal.GetTimeProvider(ctx).Now().UTC()
	for _, jb := range jobs.Jobs {
		q := queues[jb.Name]
		// fmt.Println("ASDFSADF", q.Version)
		// fmt.Printf("Enqueue: %+v\n", jb.Data)

		jb.EnqueuedAt = now
		jb.Version = resource.NewVersion(1)
		jb.QueueVersion = q.Version
		jb.Status = resource.NewStatus(resource.StatusQueued)
	}

	ids, err := be.writeJobs(ctx, jobs.Jobs)
	if err != nil {
		return nil, err
	}
	// if a job is v1, the *stored* data wont have an id. the first time we
	// update the job, the id must be stored.
	for i, jb := range jobs.Jobs {
		jb.ID = ids[i]
	}

	// now that we have ids, index. yes this isn't atomic, but it's just an index
	_, err = be.rds.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, jb := range jobs.Jobs {
			q := queues[jb.Name]
			if err := be.indexJob(ctx, pipe, q.Name, q.Labels, jb, nil); err != nil {
				return err
			}
			if err := be.checkpointJob(ctx, pipe, jb.ID, jb.ID); err != nil {
				return err
			}
		}
		return nil
	})
	return jobs, nil
}

func (be *RedisBackend) DequeueJobs(ctx context.Context, num int, opts *resource.JobListParams) (*resource.Jobs, error) {
	if opts == nil {
		opts = &resource.JobListParams{}
	}
	opts.Statuses = []*resource.Status{resource.NewStatus(resource.StatusQueued), resource.NewStatus(resource.StatusFailed)}

	ids, err := be.indexLookupJob(ctx, int64(num), opts)
	if err != nil {
		return nil, err
	}

	pendingJobs, err := be.readJobs(ctx, ids)
	if err != nil {
		return nil, err
	}

	res, err := be.queuesForJobs(ctx, pendingJobs)
	if err != nil {
		return nil, err
	}
	queues := res.ToMap()

	now := internal.GetTimeProvider(ctx).Now().UTC()
	var jobs []*resource.Job
	pendingMap := make(map[string]*resource.Job)
	for _, pjb := range pendingJobs {
		pendingMap[pjb.ID] = pjb
		// fmt.Printf("DequeueJobs pendingJob: %+v\n", pjb)
		if pjb.Data != nil && len(pjb.Data.Claims) > 0 {
			queue := queues[pjb.Name]
			match := pjb.Data.Claims.Match(opts.Claims)
			expired := queue.ClaimExpired(pjb, now)
			fmt.Println("match", match, "expired", expired)
			if !expired && !match {
				continue
			}
		}
		jb := pjb.Copy()
		jb.Version.Inc()
		status := resource.NewStatus(resource.StatusRunning)
		jb.Status = status

		jb.Results = []*resource.JobResult{
			{
				StartedAt: now,
				Status:    status,
				// TODO Attempt:
			},
		}

		jobs = append(jobs, jb)
	}

	// write new jobs
	newIds, err := be.writeJobs(ctx, jobs)
	if err != nil {
		return nil, err
	}

	// update indexes
	_, err = be.rds.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i, jb := range jobs {
			// fmt.Printf("UHUHUHUH %+v\n", jb)
			prev := pendingMap[jb.ID]
			if err := be.indexJob(ctx, pipe, jb.Name, nil, jb, prev); err != nil {
				return err
			}
			if err := be.checkpointJob(ctx, pipe, jb.ID, newIds[i]); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// if len(jobs) > 0 {
	// 	fmt.Printf("JOOOOB: %+v\n", jobs[0])
	// }
	return &resource.Jobs{Jobs: jobs}, nil
}

func (be *RedisBackend) AckJobs(ctx context.Context, req *resource.Acks) error {
	ids, err := be.lookupCheckpoints(ctx, req.JobIDs())
	if err != nil {
		return err
	}

	runningJobs, err := be.readJobs(ctx, ids)
	if err != nil {
		return err
	}

	now := internal.GetTimeProvider(ctx).Now().UTC()
	jobs := make([]*resource.Job, len(runningJobs))
	for i, rjb := range runningJobs {
		ack := req.Acks[i]

		jb := rjb.Copy()
		// if it's already complete, we had multiple runs, don't set one to
		// failed if it succeeded in a concurrent run
		shouldInc := false
		// don't update jobs that aren't running
		if *jb.Status != resource.StatusRunning {
			continue
		}

		// don't change from complete, in case there are concurrent jobs for the same id
		if *jb.Status != resource.StatusComplete {
			if jb.Status != ack.Status {
				shouldInc = true
				jb.Status = ack.Status
			}
		}

		shouldInc = true
		res := jb.LastResult()
		res.CompletedAt = now
		// res.Attempt =
		if ack.Data != nil {
			res.Data = ack.Data
		}
		res.Status = ack.Status
		res.Error = ack.Error

		if shouldInc {
			jb.Version.Inc()
		}
		jobs[i] = jb
	}

	// write new jobs
	newIds, err := be.writeJobs(ctx, jobs)
	if err != nil {
		return err
	}
	// update indexes
	_, err = be.rds.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i, jb := range jobs {
			prev := runningJobs[i]
			if err := be.indexJob(ctx, pipe, jb.Name, nil, jb, prev); err != nil {
				return err
			}
			if err := be.checkpointJob(ctx, pipe, jb.ID, newIds[i]); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (be *RedisBackend) GetSetJobKeys(ctx context.Context, keys []string) (bool, error) {
	return false, nil
}

func (be *RedisBackend) DeleteJobKeys(ctx context.Context, keys []string) error {
	return nil
}

func (be *RedisBackend) GetJobByID(ctx context.Context, idArg string) (*resource.Job, error) {
	ids, err := be.lookupCheckpoints(ctx, []string{idArg})
	if err != nil {
		return nil, err
	}
	jobs, err := be.readJobs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return jobs[0], nil
}

func (be *RedisBackend) ListJobs(ctx context.Context, limit int, opts *resource.JobListParams) (*resource.Jobs, error) {
	if limit <= 0 {
		limit = 100
	}
	ids, err := be.indexLookupJob(ctx, int64(limit), opts)
	if err != nil {
		return nil, err
	}

	jobs, err := be.readJobs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return &resource.Jobs{Jobs: jobs}, nil
}

func (be *RedisBackend) writeJobs(ctx context.Context, jobs []*resource.Job) ([]string, error) {
	cmds := make([]*redis.StringCmd, len(jobs))
	_, err := be.rds.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i, jb := range jobs {
			// fmt.Printf("writeJobs job #%d: %+v\n", i, jb.Data)
			b, err := jobv1.MarshalJob(jb)
			if err != nil {
				return err
			}
			// fmt.Printf("writeJobs job #%d: %q\n", i, string(b))

			id := jb.ID
			if id == "" {
				id = "new"
			}
			cmds[i] = pipe.XAdd(ctx, &redis.XAddArgs{
				Stream:       streamKey,
				MaxLenApprox: int64(be.cfg.MaxStreamSize),
				Values:       []string{id, string(b)},
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(jobs))
	for i, cmd := range cmds {
		id := cmd.Val()
		ids[i] = id
	}
	return ids, nil
}

func (be *RedisBackend) readJobs(ctx context.Context, ids []string) ([]*resource.Job, error) {
	jobCmds := make([]*redis.XMessageSliceCmd, len(ids))
	_, err := be.rds.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i, id := range ids {
			jobCmds[i] = pipe.XRange(ctx, streamKey, id, id)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	jobs := make([]*resource.Job, len(ids))
	for i, jobCmd := range jobCmds {
		msgs, err := jobCmd.Result()
		if err != nil {
			return nil, err
		}
		if len(msgs) != 1 {
			return nil, fmt.Errorf("backend/redis: got %d results, expected 1", len(msgs))
		}
		if len(msgs[0].Values) != 1 {
			return nil, fmt.Errorf("backend/redis: got %d keys, expected 1", len(msgs[0].Values))
		}

		vals := msgs[0].Values
		var key string
		for k := range vals {
			key = k
			break
		}
		id := msgs[0].ID
		if key != "new" {
			id = key
		}
		b := msgs[0].Values[key].(string)

		jb, err := jobv1.UnmarshalJob([]byte(b), nil)
		if err != nil {
			return nil, err
		}
		// fmt.Printf("unmarshaled job: %+v\n", jb.Data)

		if jb.ID == "" {
			jb.ID = id
		}
		jobs[i] = jb
	}

	return jobs, nil
}

func (be *RedisBackend) lookupCheckpoints(ctx context.Context, ids []string) ([]string, error) {
	res, err := be.rds.HMGet(ctx, checkpointKey, ids...).Result()
	if err != nil {
		return nil, err
	}

	resIds := make([]string, len(ids))
	for i, iid := range res {
		resIds[i] = iid.(string)
	}
	return resIds, nil
}

// TODO this needs to account for queue version
func (be *RedisBackend) queuesForJobs(ctx context.Context, jobs []*resource.Job) (*resource.Queues, error) {
	qMap := make(map[string]bool)
	for _, jb := range jobs {
		qMap[jb.Name] = true
	}
	names := make([]string, len(qMap))
	i := 0
	for name := range qMap {
		names[i] = name
		i++
	}

	return be.GetQueues(ctx, names)
}
