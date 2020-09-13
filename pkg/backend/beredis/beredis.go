// Package beredis implements the backend interface using redis, primarily
// streams.
package beredis

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v8"

	"github.com/jeffrom/job-manager/pkg/backend"
)

// RedisBackend implements backend.Interface. Indexing is eventually consistent
// right now.
// TODO probably most of this needs to be lua
type RedisBackend struct {
	cfg Config
	rds *redis.Client
}

type ProviderFunc func(b *RedisBackend) *RedisBackend

func New(providers ...ProviderFunc) *RedisBackend {
	be := &RedisBackend{
		cfg: defaultConfig,
	}
	for _, pr := range providers {
		be = pr(be)
	}

	be.rds = redis.NewClient(be.cfg.Redis)
	return be
}

func WithConfig(cfg Config) ProviderFunc {
	return func(b *RedisBackend) *RedisBackend {
		b.cfg = cfg
		return b
	}
}

func (be *RedisBackend) Ping(ctx context.Context) error {
	return be.rds.Ping(ctx).Err()
}

func (be *RedisBackend) Reset(ctx context.Context) error {
	if !be.cfg.TestMode {
		return backend.ErrNotAuthorized
	}
	keys, err := be.rds.Keys(ctx, "mjob:*").Result()
	if err != nil {
		return err
	}
	keys = append(keys, streamKey, queueListKey)
	return be.rds.Del(ctx, keys...).Err()
}

func redisKey(parts ...string) string { return strings.Join(parts, ":") }
