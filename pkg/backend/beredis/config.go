package beredis

import (
	"github.com/go-redis/redis/v8"

	"github.com/jeffrom/job-manager/pkg/backend"
)

type Config struct {
	backend.Config
	Redis *redis.Options `json:"-"`
}

var defaultConfig = Config{
	Config: backend.DefaultConfig,
	Redis: &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	},
}
