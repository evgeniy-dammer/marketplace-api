package repository

import (
	"fmt"
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// database table names
const (
	userTable         = "users"
	roleTable         = "roles"
	statusTable       = "statuses"
	userRoleTable     = "users_roles"
	organisationTable = "organisations"
	categoryTable     = "categories"
	itemTable         = "items"
	categoryItemTable = "categories_items"
)

// NewPostgresDB create connection to database
func NewPostgresDB(cfg model.DbConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.DbName, cfg.Password, cfg.SSLMode,
		),
	)

	if err != nil {
		return nil, err
	}

	// verify connection
	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
