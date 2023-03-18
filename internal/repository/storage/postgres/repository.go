package postgres

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	genSQL      squirrel.StatementBuilderType
	database    *sqlx.DB
	options     Options
	isTracingOn bool
}

type Options struct {
	Timeout time.Duration
}

func New(database *sqlx.DB, options Options, isTracingOn bool) *Repository {
	return &Repository{
		genSQL:      squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		database:    database,
		options:     options,
		isTracingOn: isTracingOn,
	}
}
