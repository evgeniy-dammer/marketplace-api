package db

import (
	"context"
	"fmt"
	"os"

	"github.com/evgeniy-dammer/emenu-api/pkg/common/models"
	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

// Connect establishes a connection with a PostgreSQL database
func Connect(c *models.DbConfiguration) {
	var err error

	databaseUrl := "postgres://" + c.DbUser + ":" + c.DbPass + "@" + c.DbHost + ":" + c.DbPort + "/" + c.DbName

	DB, err = pgx.Connect(context.Background(), databaseUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}
