package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Database connection pool
var DBPool *pgxpool.Pool

// ConnectToDB initializes the database connection using the pgx connection pool.
func ConnectToDB() error {
    var err error

    // Retrieve the DATABASE_URL from the environment variables
    databaseURL := os.Getenv("DATABASE_URL")
    if databaseURL == "" {
        return fmt.Errorf("DATABASE_URL is not set")
    }

    // Connect to the database
    DBPool, err = pgxpool.Connect(context.Background(), databaseURL)
    if err != nil {
        return fmt.Errorf("Unable to connect to database: %v", err)
    }

    fmt.Println("Connected to the database successfully")
    return nil
}

// CloseDB closes the database connection pool.
func CloseDB() {
    DBPool.Close()
}
