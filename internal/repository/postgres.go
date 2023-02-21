package repository

import (
	"fmt"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// database table names.
const (
	userTable         = "users"
	roleTable         = "roles"
	statusTable       = "users_statuses"
	userRoleTable     = "users_roles"
	organizationTable = "organizations"
	categoryTable     = "categories"
	itemTable         = "items"
	tableTable        = "tables"
	orderTable        = "orders"
	orderItemTable    = "orders_items"
	imageTable        = "images"

	// categoryItemTable = "categories_items"
)

// NewPostgresDB create connection to database.
func NewPostgresDB(cfg model.DBConfig) (*sqlx.DB, error) {
	database, err := sqlx.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode,
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open connection to database")
	}

	// verify connection
	err = database.Ping()

	if err != nil {
		return nil, errors.Wrap(err, "unable to ping connection to database")
	}

	return database, nil
}
