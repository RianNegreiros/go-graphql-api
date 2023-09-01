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

func (ur *UserRepo) Create(ctx context.Context, user models.User) (models.User, error) {
	tx, err := ur.DB.Pool.Begin(ctx)
	if err != nil {
		return models.User{}, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	user, err = createUser(ctx, tx, user)
	if err != nil {
		return models.User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.User{}, fmt.Errorf("error commiting: %v", err)
	}

	return user, nil
}

func createUser(ctx context.Context, tx pgx.Tx, user models.User) (models.User, error) {
	query := `INSERT INTO users (email, username, password) VALUES ($1, $2, $3) RETURNING *;`

	u := models.User{}

	if err := pgxscan.Get(ctx, tx, &u, query, user.Email, user.Username, user.Password); err != nil {
		return models.User{}, fmt.Errorf("error insert: %v", err)
	}

	return u, nil
}

func (ur *UserRepo) GetByUsername(ctx context.Context, username string) (models.User, error) {
	query := `SELECT * FROM users WHERE username = $1 LIMIT 1;`

	u := models.User{}

	if err := pgxscan.Get(ctx, ur.DB.Pool, &u, query, username); err != nil {
		if pgxscan.NotFound(err) {
			return models.User{}, models.ErrNotFound
		}

		return models.User{}, fmt.Errorf("error select: %v", err)
	}

	return u, nil
}

func (ur *UserRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	query := `SELECT * FROM users WHERE email = $1 LIMIT 1;`

	u := models.User{}

	if err := pgxscan.Get(ctx, ur.DB.Pool, &u, query, email); err != nil {
		if pgxscan.NotFound(err) {
			return models.User{}, models.ErrNotFound
		}

		return models.User{}, fmt.Errorf("error select: %v", err)
	}

	return u, nil
}

func (ur *UserRepo) GetByID(ctx context.Context, id string) (models.User, error) {
	query := `SELECT * FROM users WHERE id = $1 LIMIT 1;`

	u := models.User{}

	if err := pgxscan.Get(ctx, ur.DB.Pool, &u, query, id); err != nil {
		if pgxscan.NotFound(err) {
			return models.User{}, models.ErrNotFound
		}

		return models.User{}, fmt.Errorf("error select: %v", err)
	}

	return u, nil
}

func (ur *UserRepo) GetByIds(ctx context.Context, ids []string) ([]models.User, error) {
	return getUsersByIds(ctx, ur.DB.Pool, ids)
}

func getUsersByIds(ctx context.Context, q pgxscan.Querier, ids []string) ([]models.User, error) {
	query := `SELECT * FROM users WHERE id = ANY($1);`

	var uu []models.User

	if err := pgxscan.Select(ctx, q, &uu, query, ids); err != nil {
		return nil, fmt.Errorf("error get users by ids: %+v", err)
	}

	return uu, nil
}
