package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Printf("Unable to connect to database: %v", err)
		return nil, err
	}

	log.Printf("Successfully connected to database")

	return pool, nil
}
