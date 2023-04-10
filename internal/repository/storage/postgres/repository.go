package postgres

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	genSQL         squirrel.StatementBuilderType
	databaseMaster *sqlx.DB
	databaseSlave  *sqlx.DB
	options        Options
	isTracingOn    bool
}

type Options struct {
	Timeout time.Duration
}

func New(databaseMaster *sqlx.DB, databaseSlave *sqlx.DB, options Options, isTracingOn bool) *Repository {
	return &Repository{
		genSQL:         squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		databaseMaster: databaseMaster,
		databaseSlave:  databaseSlave,
		options:        options,
		isTracingOn:    isTracingOn,
	}
}
