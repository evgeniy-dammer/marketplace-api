package db

import (
	"context"

	"github.com/evgeniy-dammer/emenu-api/pkg/common/models"
	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

// Connect establishes a connection with a PostgreSQL database
func Connect(c *models.Configuration) error {
	var err error

	databaseUrl := "postgres://" + c.DbUser + ":" + c.DbPass + "@" + c.DbHost + ":" + c.DbPort + "/" + c.DbName

	DB, err = pgx.Connect(context.Background(), databaseUrl)

	return err
}
