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

// const queueKey = "mjob:q"

func (be *RedisBackend) EnqueueJobs(ctx context.Context, jobs *resource.Jobs) (*resource.Jobs, error) {
	// first check everything has a queue
	res, err := be.queuesForJobs(ctx, jobs)
	if err != nil {
		return nil, err
	}
	queues := res.ToMap()

	now := internal.GetTimeProvider(ctx).Now().UTC()
	for _, jb := range jobs.Jobs {
		q := queues[jb.Name]
		fmt.Println("ASDFSADF", q.Version)

		jb.EnqueuedAt = now
		jb.Version = resource.NewVersion(1)
		jb.QueueVersion = q.Version
		jb.Status = resource.StatusQueued
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
			if err := be.indexJob(ctx, pipe, q.ID, jb, nil); err != nil {
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
	opts.Statuses = []resource.Status{resource.StatusQueued, resource.StatusFailed}

	ids, err := be.indexLookup(ctx, int64(num), opts)
	if err != nil {
		return nil, err
	}

	pendingJobs, err := be.readJobs(ctx, ids)
	if err != nil {
		return nil, err
	}

	jobs := make([]*resource.Job, len(pendingJobs))
	for i, pjb := range pendingJobs {
		jb := pjb.Copy()
		jb.Version.Inc()
		jb.Status = resource.StatusRunning
		jobs[i] = jb
	}

	// write new jobs
	newIds, err := be.writeJobs(ctx, jobs)
	if err != nil {
		return nil, err
	}

	// update indexes
	_, err = be.rds.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i, jb := range jobs {
			fmt.Printf("UHUHUHUH %+v\n", jb)
			prev := pendingJobs[i]
			if err := be.indexJob(ctx, pipe, jb.Name, jb, prev); err != nil {
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
	if len(jobs) > 0 {
		fmt.Printf("JOOOOB: %+v\n", jobs[0])
	}
	return &resource.Jobs{Jobs: jobs}, nil
}

func (be *RedisBackend) AckJobs(ctx context.Context, req *resource.Acks) error {
	ids, err := be.lookupCheckpoints(ctx, req.IDs())
	if err != nil {
		return err
	}

	runningJobs, err := be.readJobs(ctx, ids)
	if err != nil {
		return err
	}

	jobs := make([]*resource.Job, len(runningJobs))
	for i, rjb := range runningJobs {
		ack := req.Acks[i]

		jb := rjb.Copy()
		// if it's already complete, we had multiple runs, don't set one to
		// failed if it succeeded in a concurrent run
		shouldInc := false
		if jb.Status != resource.StatusComplete {
			if jb.Status != ack.Status {
				shouldInc = true
				jb.Status = ack.Status
			}
		}

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
			if err := be.indexJob(ctx, pipe, jb.Name, jb, prev); err != nil {
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

func (be *RedisBackend) GetJobByID(ctx context.Context, id string) (*resource.Job, error) {
	return nil, nil
}

func (be *RedisBackend) ListJobs(ctx context.Context, limit int, opts *resource.JobListParams) (*resource.Jobs, error) {
	if limit <= 0 {
		limit = 100
	}
	ids, err := be.indexLookup(ctx, int64(limit), opts)
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
			b, err := jobv1.MarshalJob(jb)
			if err != nil {
				return err
			}

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

func (be *RedisBackend) queuesForJobs(ctx context.Context, jobs *resource.Jobs) (*resource.Queues, error) {
	qMap := make(map[string]bool)
	for _, jb := range jobs.Jobs {
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
