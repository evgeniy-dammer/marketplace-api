package postgres

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	options     Options
	database    *sqlx.DB
	isTracingOn bool
}

type Options struct {
	Timeout time.Duration
}

func New(database *sqlx.DB, options Options, isTracingOn bool) *Repository {
	repository := &Repository{
		database:    database,
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
