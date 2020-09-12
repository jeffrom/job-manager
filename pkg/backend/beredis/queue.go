package beredis

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"

	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/resource"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
)

const queueListKey = "mjob:queues"

func queueKey(name string) string {
	return "mjob:queues:" + name
}

func (be *RedisBackend) GetQueue(ctx context.Context, job string) (*resource.Queue, error) {
	return nil, nil
}

func (be *RedisBackend) SaveQueue(ctx context.Context, queue *resource.Queue) (*resource.Queue, error) {
	if queue == nil || queue.ID == "" {
		return nil, backend.ErrInvalidResource
	}

	key := queueKey(queue.ID)

	// TODO this could be TxPipeline? use pop instead of trim?
	err := be.rds.Watch(ctx, func(tx *redis.Tx) error {
		queueV := queue.Version
		if queueV == nil {
			queueV = resource.NewVersion(0)
		}
		prevb, err := tx.LIndex(ctx, key, -1).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return err
		}

		prev, err := jobv1.UnmarshalQueue([]byte(prevb), nil)
		if err != nil {
			return err
		}

		if prev != nil && queue.Equals(prev) {
			return nil
		}
		queueV.Inc()
		queue.Version = queueV

		l, err := tx.LLen(ctx, key).Result()
		if err != nil {
			return err
		}
		if limit := be.cfg.HistoryLimit; limit > 0 && l >= int64(limit) {
			delta := l - int64(limit)
			if err := tx.LTrim(ctx, key, delta+1, -1).Err(); err != nil {
				return err
			}
		}

		b, err := jobv1.MarshalQueue(queue)
		if err != nil {
			return err
		}
		if err := tx.RPush(ctx, key, string(b)).Err(); err != nil {
			return err
		}

		// update the queue list
		return nil
	}, key, queueListKey)
	if err != nil {
		return nil, err
	}
	return queue, nil
}

func (be *RedisBackend) ListQueues(ctx context.Context, opts *resource.QueueListParams) (*resource.Queues, error) {
	return nil, nil
}
