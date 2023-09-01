package postgres

import (
	"context"
	"fmt"

	"github.com/RianNegreiros/go-graphql-api/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type UserRepo struct {
	DB *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (r *UserRepo) Create(ctx context.Context, user models.User) (models.User, error) {
	tx, err := r.DB.Pool.Begin(ctx)
	if err != nil {
		return models.User{}, fmt.Errorf("error starting transaction: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("error rolling back transaction: %v", err)
		}
	}(tx, ctx)

	user, err = createUser(ctx, tx, user)
	if err != nil {
		return models.User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.User{}, fmt.Errorf("error committing transaction: %w", err)
	}

	return user, nil
}

func createUser(ctx context.Context, tx pgx.Tx, user models.User) (models.User, error) {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING *;;`

	u := models.User{}

	if err := pgxscan.Get(ctx, tx, &u, query, user.Username, user.Email, user.Password); err != nil {
		return models.User{}, fmt.Errorf("error creating user: %w", err)
	}

	return u, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (models.User, error) {
	query := `SELECT * FROM users WHERE username = $1 LIMIT 1;`

	u := models.User{}

	if err := pgxscan.Get(ctx, r.DB.Pool, &u, query, username); err != nil {
		if pgxscan.NotFound(err) {
			return models.User{}, models.ErrNotFound
		}

		return models.User{}, fmt.Errorf("error getting user by username: %w", err)
	}

	return u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	query := `SELECT * FROM users WHERE email = $1 LIMIT 1;`

	u := models.User{}

	if err := pgxscan.Get(ctx, r.DB.Pool, &u, query, email); err != nil {
		if pgxscan.NotFound(err) {
			return models.User{}, models.ErrNotFound
		}

		return models.User{}, fmt.Errorf("error getting user by username: %w", err)
	}

	return u, nil
}
