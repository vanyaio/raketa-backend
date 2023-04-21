package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "postgres", "postgres", "raketadb",
	)

	dbpool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}

	if err := dbpool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	return dbpool, nil
}
