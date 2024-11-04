package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

var DB *pgxpool.Pool

func ConnectDB() {
	dbURL := os.Getenv("COCKROACHDB_URL")
	var err error
	DB, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	fmt.Println("Connected to CockroachDB!")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
