package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect() *sql.DB {
	connectionString := os.Getenv("POSTGRES_URL")

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open database handle: %v\n", err)
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return db
}
