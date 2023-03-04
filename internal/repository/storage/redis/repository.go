package redis

import (
	"github.com/go-redis/cache/v8"
)

type Repository struct {
	cache *cache.Cache
}

func New(cache *cache.Cache) *Repository {
	return &Repository{cache: cache}
}
