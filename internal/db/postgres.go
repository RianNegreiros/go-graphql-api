package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/RianNegreiros/go-graphql-api/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"runtime"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type DB struct {
	Pool   *pgxpool.Pool
	config *config.Config
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
		Pool:   pool,
		config: config,
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

func (db *DB) Migrate() error {
	_, _, _, _ = runtime.Caller(0)

	migrationPath := fmt.Sprintf("file://internal/db/migrations")

	m, err := migrate.New(migrationPath, db.config.Database.URL)
	if err != nil {
		return fmt.Errorf("error creating migration instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error running migrations: %w", err)
	}

	log.Println("migrations ran successfully")

	return nil
}
