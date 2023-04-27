package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/pgxpool"
	"github.com/vanyaio/raketa-backend/config"
)

func NewPool(ctx context.Context, config *config.Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Postgres.PostgresqlHost, config.Postgres.PostgresqlUser, 
		config.Postgres.PostgresqlPassword, config.Postgres.PostgresqlDbname,
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
