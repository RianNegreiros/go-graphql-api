package postgres

import (
	"context"
	"fmt"

	"github.com/RianNegreiros/go-graphql-api/internal/user"
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

func (ur *UserRepo) Create(ctx context.Context, userModel user.UserModel) (user.UserModel, error) {
	tx, err := ur.DB.Pool.Begin(ctx)
	if err != nil {
		return user.UserModel{}, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	userModel, err = createUser(ctx, tx, userModel)
	if err != nil {
		return user.UserModel{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return user.UserModel{}, fmt.Errorf("error commiting: %v", err)
	}

	return userModel, nil
}

func createUser(ctx context.Context, tx pgx.Tx, userModel user.UserModel) (user.UserModel, error) {
	query := `INSERT INTO users (email, username, password) VALUES ($1, $2, $3) RETURNING *;`

	u := user.UserModel{}

	if err := pgxscan.Get(ctx, tx, &u, query, userModel.Email, userModel.Username, userModel.Password); err != nil {
		return user.UserModel{}, fmt.Errorf("error insert: %v", err)
	}

	return u, nil
}

func (ur *UserRepo) GetByUsername(ctx context.Context, username string) (user.UserModel, error) {
	query := `SELECT * FROM users WHERE username = $1 LIMIT 1;`

	u := user.UserModel{}

	if err := pgxscan.Get(ctx, ur.DB.Pool, &u, query, username); err != nil {
		if pgxscan.NotFound(err) {
			return user.UserModel{}, user.ErrNotFound
		}

		return user.UserModel{}, fmt.Errorf("error select: %v", err)
	}

	return u, nil
}

func (ur *UserRepo) GetByEmail(ctx context.Context, email string) (user.UserModel, error) {
	query := `SELECT * FROM users WHERE email = $1 LIMIT 1;`

	u := user.UserModel{}

	if err := pgxscan.Get(ctx, ur.DB.Pool, &u, query, email); err != nil {
		if pgxscan.NotFound(err) {
			return user.UserModel{}, user.ErrNotFound
		}

		return user.UserModel{}, fmt.Errorf("error select: %v", err)
	}

	return u, nil
}

func (ur *UserRepo) GetByID(ctx context.Context, id string) (user.UserModel, error) {
	query := `SELECT * FROM users WHERE id = $1 LIMIT 1;`

	u := user.UserModel{}

	if err := pgxscan.Get(ctx, ur.DB.Pool, &u, query, id); err != nil {
		if pgxscan.NotFound(err) {
			return user.UserModel{}, user.ErrNotFound
		}

		return user.UserModel{}, fmt.Errorf("error select: %v", err)
	}

	return u, nil
}

func (ur *UserRepo) GetByIds(ctx context.Context, ids []string) ([]user.UserModel, error) {
	return getUsersByIds(ctx, ur.DB.Pool, ids)
}

func getUsersByIds(ctx context.Context, q pgxscan.Querier, ids []string) ([]user.UserModel, error) {
	query := `SELECT * FROM users WHERE id = ANY($1);`

	var uu []user.UserModel

	if err := pgxscan.Select(ctx, q, &uu, query, ids); err != nil {
		return nil, fmt.Errorf("error get users by ids: %+v", err)
	}

	return uu, nil
}
