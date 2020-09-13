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
	now := internal.GetTimeProvider(ctx).Now().UTC()
	// first check everything has a queue
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

	res, err := be.GetQueues(ctx, names)
	if err != nil {
		return nil, err
	}

	queues := make(map[string]*resource.Queue)
	for _, q := range res.Queues {
		queues[q.ID] = q
	}

	cmds := make([]*redis.StringCmd, len(jobs.Jobs))
	_, err = be.rds.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i, jb := range jobs.Jobs {
			q := queues[jb.Name]
			if jb.EnqueuedAt.IsZero() {
				jb.EnqueuedAt = now
			}
			if v := jb.Version; v == nil {
				jb.Version = resource.NewVersion(1)
			}
			if v := jb.QueueVersion; v == nil {
				jb.QueueVersion = q.Version
			}
			jb.Status = resource.StatusQueued

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

	for i, cmd := range cmds {
		id := cmd.Val()
		// if a job is v1, the *stored* data wont have an id. the first time we
		// update the job, the id must be stored.
		jobs.Jobs[i].ID = id
	}

	// now that we have ids, index. yes this isn't atomic, but it's just an index
	_, err = be.rds.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, jb := range jobs.Jobs {
			q := queues[jb.Name]
			if err := be.indexJob(ctx, pipe, q.ID, jb, nil); err != nil {
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

	jobCmds := make([]*redis.XMessageSliceCmd, len(ids))
	_, err = be.rds.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i, id := range ids {
			jobCmds[i] = pipe.XRange(ctx, streamKey, id, id)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var jobs []*resource.Job
	for _, jobCmd := range jobCmds {
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
		jb.Version.Inc()
		jb.Status = resource.StatusRunning
		jobs = append(jobs, jb)
	}
	return &resource.Jobs{Jobs: jobs}, nil
}

func (be *RedisBackend) AckJobs(ctx context.Context, results *resource.Acks) error {

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

func (be *RedisBackend) ListJobs(ctx context.Context, opts *resource.JobListParams) (*resource.Jobs, error) {
	return nil, nil
}
