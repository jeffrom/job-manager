package beredis

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"

	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/internal"
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
	now := internal.GetTimeProvider(ctx).Now()
	key := queueKey(queue.ID)

	// TODO this could be TxPipeline? use pop instead of trim?
	err := be.rds.Watch(ctx, func(tx *redis.Tx) error {
		queueV := queue.Version
		if queueV == nil {
			return backend.ErrInvalidResource
		}
		prevb, err := tx.LIndex(ctx, key, -1).Result()
		keyIsNil := errors.Is(err, redis.Nil)
		if err != nil && !keyIsNil {
			return err
		}

		var prev *resource.Queue
		if keyIsNil {
			queue.CreatedAt = now
			queue.UpdatedAt = now
		} else {
			prev, err = jobv1.UnmarshalQueue([]byte(prevb), nil)
			if err != nil {
				return err
			}
		}

		if prev != nil {
			queue.CreatedAt = prev.CreatedAt
			queue.UpdatedAt = prev.UpdatedAt
		}
		// fmt.Printf("---\nprev: %+v\n", prev)
		// fmt.Printf("curr: %+v\n", queue)
		if prev != nil && queue.Equals(prev) {
			return nil
		} else if prev != nil && !prev.Version.Equals(queueV) {
			return &backend.VersionConflictError{
				Resource:   "queue",
				ResourceID: queue.ID,
				Prev:       prev.Version.String(),
				Curr:       queue.Version.String(),
			}
		}

		if !keyIsNil {
			queue.UpdatedAt = now
			queueV.Inc()
		}

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
