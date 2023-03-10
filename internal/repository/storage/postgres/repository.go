package postgres

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	options  Options
	database *sqlx.DB
}

type Options struct {
	Timeout time.Duration
}

func New(database *sqlx.DB, options Options) *Repository {
	repository := &Repository{
		database: database,
	}

	repository.SetOptions(options)

	return repository
}

func (r *Repository) SetOptions(options Options) {
	if r.options != options {
		r.options = options
	}
}
