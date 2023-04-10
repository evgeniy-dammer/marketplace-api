package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	client      *redis.Client
	options     Options
	isTracingOn bool
}

type Options struct {
	Timeout time.Duration
	TTL     time.Duration
}

func New(client *redis.Client, options Options, isTracingOn bool) *Repository {
	repository := &Repository{
		client:      client,
		isTracingOn: isTracingOn,
	}

	repository.SetOptions(options)

	return repository
}

func (r *Repository) SetOptions(options Options) {
	if r.options != options {
		r.options = options
	}
}
