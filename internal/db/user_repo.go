package db

import (
	"context"
	"fmt"
	"github.com/RianNegreiros/go-graphql-api/internal"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type UserRepo struct {
	DB *DB
}

func (r *UserRepo) Create(ctx context.Context, user internal.User) (internal.User, error) {
	tx, err := r.DB.Pool.Begin(ctx)
	if err != nil {
		return internal.User{}, fmt.Errorf("error starting transaction: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("error rolling back transaction: %v", err)
		}
	}(tx, ctx)

	user, err = createUser(ctx, tx, user)
	if err != nil {
		return internal.User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return internal.User{}, fmt.Errorf("error committing transaction: %w", err)
	}

	return user, nil
}

func createUser(ctx context.Context, tx pgx.Tx, user internal.User) (internal.User, error) {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING *;;`

	u := internal.User{}

	if err := pgxscan.Get(ctx, tx, &u, query); err != nil {
		return internal.User{}, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (internal.User, error) {
	query := `SELECT * FROM users WHERE username = $1 LIMIT 1;`

	u := internal.User{}

	if err := pgxscan.Get(ctx, r.DB.Pool, &u, query, username); err != nil {
		if pgxscan.NotFound(err) {
			return internal.User{}, internal.ErrNotFound
		}

		return internal.User{}, fmt.Errorf("error getting user by username: %w", err)
	}

	return u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (internal.User, error) {
	query := `SELECT * FROM users WHERE email = $1 LIMIT 1;`

	u := internal.User{}

	if err := pgxscan.Get(ctx, r.DB.Pool, &u, query, email); err != nil {
		if pgxscan.NotFound(err) {
			return internal.User{}, internal.ErrNotFound
		}

		return internal.User{}, fmt.Errorf("error getting user by username: %w", err)
	}

	return u, nil
}
