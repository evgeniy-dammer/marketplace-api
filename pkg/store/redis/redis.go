package redis

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// NewRedisCache create connection to redis.
func NewRedisCache(cfg redis.Options) (*cache.Cache, error) {
	client := redis.NewClient(&cfg)

	if err := client.Ping(context.Background()); err != nil {
		return nil, errors.Wrap(err.Err(), "unable to connect to redis")
	}

	redisCache := cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return redisCache, nil
}
