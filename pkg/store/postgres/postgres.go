package postgres

import (
	"fmt"
	"net"

	"github.com/casbin/casbin-pg-adapter"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/pressly/goose"
)

// DBConfig is a database config.
type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB create connection to database.
func NewPostgresDB(cfg DBConfig) (*sqlx.DB, *pgadapter.Adapter, error) {
	database, err := sqlx.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode,
		),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to open connection to database")
	}

	err = migrationsUp(database)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to up migrations")
	}

	adapter, err := pgadapter.NewAdapter(
		fmt.Sprintf(
			"postgresql://%s:%s@%s/%s?sslmode=%s",
			cfg.Username, cfg.Password, net.JoinHostPort(cfg.Host, cfg.Port), cfg.DBName, cfg.SSLMode,
		),
		cfg.DBName,
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to open adapter connection to database")
	}

	err = database.Ping()
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to ping connection to database")
	}

	return database, adapter, nil
}

// migrationsUp up database migrations
func migrationsUp(database *sqlx.DB) error {
	if err := goose.Run("up", database.DB, "./schema"); err != nil {
		return errors.Wrap(err, "error occurred while running goose up command")
	}

	return nil
}
