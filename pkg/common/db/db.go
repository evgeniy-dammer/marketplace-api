package db

import (
	"context"
	"fmt"
	"os"

	"github.com/evgeniy-dammer/emenu-api/pkg/common/config"
	"github.com/jackc/pgx/v5"
)

func Connect(c *config.Configuration) *pgx.Conn {
	databaseUrl := "postgres://" + c.DBUser + ":" + c.DBPass + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName

	connection, err := pgx.Connect(context.Background(), databaseUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database connected!")

	return connection
}
