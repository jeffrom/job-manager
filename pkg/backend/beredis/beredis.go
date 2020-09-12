// Package beredis implements the backend interface using redis, primarily
// streams.
package beredis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

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
