package beredis

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/jeffrom/job-manager/pkg/resource"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
)

const streamKey = "mjob:jobs"

func (be *RedisBackend) EnqueueJobs(ctx context.Context, jobs *resource.Jobs) (*resource.Jobs, error) {
	cmds := make([]*redis.StringCmd, len(jobs.Jobs))
	_, err := be.rds.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i, jb := range jobs.Jobs {
			b, err := jobv1.MarshalJob(jb)
			if err != nil {
				return err
			}

			cmds[i] = pipe.XAdd(ctx, &redis.XAddArgs{
				Stream:       streamKey,
				MaxLenApprox: int64(be.cfg.MaxStreamSize),
				Values:       []string{"job", string(b)},
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	for i, cmd := range cmds {
		id := cmd.Val()
		jobs.Jobs[i].ID = id
	}
	return jobs, nil
}

func (be *RedisBackend) DequeueJobs(ctx context.Context, num int, opts *resource.JobListParams) (*resource.Jobs, error) {

	return nil, nil
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
