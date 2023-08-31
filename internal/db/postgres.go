package db

import (
	"context"
	"github.com/RianNegreiros/go-graphql-api/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, config *config.Config) *DB {
	dbConfig, err := pgxpool.ParseConfig(config.Database.URL)
	if err != nil {
		log.Fatalf("error loading db config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Fatalf("error connecting to postgres: %v", err)
	}

	db := &DB{
		Pool: pool,
	}

	db.Ping(ctx)

	return db
}

func (db *DB) Ping(ctx context.Context) {
	if err := db.Pool.Ping(ctx); err != nil {
		log.Fatalf("error pinging db: %v", err)
	}

	log.Println("connected to postgres")
}

func (db *DB) Close() {
	db.Pool.Close()
}
