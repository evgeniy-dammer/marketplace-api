package db

import (
	"context"
	"fmt"
	"os"

	"github.com/evgeniy-dammer/emenu-api/pkg/common/models"
	"github.com/jackc/pgx/v5"
)

// Connect establishes a connection with a PostgreSQL database
func Connect(c *models.DbConfiguration) *pgx.Conn {
	databaseUrl := "postgres://" + c.DbUser + ":" + c.DbPass + "@" + c.DbHost + ":" + c.DbPort + "/" + c.DbName

	connection, err := pgx.Connect(context.Background(), databaseUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database connected!")

	return connection
}
