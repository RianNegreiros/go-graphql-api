package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path"
	"runtime"

	"github.com/RianNegreiros/go-graphql-api/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var migrationPath string

type DB struct {
	Pool   *pgxpool.Pool
	config *config.Config
}

func New(ctx context.Context, config *config.Config) *DB {
	dbConfig, err := pgxpool.ParseConfig(config.Database.URL)
	if err != nil {
		log.Fatalf("error loading postgres config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Fatalf("error connecting to postgres: %v", err)
	}

	if config.Env.BuildEnv == "docker" {
		migrationPath = "file://./migrations"
	}

	_, b, _, _ := runtime.Caller(0)

	migrationPath = fmt.Sprintf("file://%s/migrations", path.Dir(b))

	db := &DB{
		Pool:   pool,
		config: config,
	}

	db.Ping(ctx)

	return db
}

func (db *DB) Ping(ctx context.Context) {
	if err := db.Pool.Ping(ctx); err != nil {
		log.Fatalf("error pinging postgres: %v", err)
	}

	log.Println("connected to postgres")
}

func (db *DB) Close() {
	db.Pool.Close()
}

func (db *DB) Migrate() error {
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

func (db *DB) Drop() error {
	m, err := migrate.New(migrationPath, db.config.Database.URL)
	if err != nil {
		return fmt.Errorf("error creating migration instance: %w", err)
	}

	if err := m.Drop(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error dropping migrations: %w", err)
	}

	log.Println("migrations dropped successfully")

	return nil
}

func (db *DB) Truncate(ctx context.Context) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM users;`)
	if err != nil {
		return fmt.Errorf("error truncating users table: %w", err)
	}

	return nil
}
