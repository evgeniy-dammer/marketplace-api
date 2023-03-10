package redis

import (
	"time"

	"github.com/go-redis/cache/v8"
)

type Repository struct {
	options Options
	cache   *cache.Cache
}

type Options struct {
	Timeout time.Duration
}

func New(cache *cache.Cache, options Options) *Repository {
	repository := &Repository{
		cache: cache,
	}

	repository.SetOptions(options)

	return repository
}

func (r *Repository) SetOptions(options Options) {
	if r.options != options {
		r.options = options
	}
}
