package user

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUsernameTaken = errors.New("username already taken")
	ErrEmailTaken    = errors.New("email already taken")
)

type UserService interface {
	GetByID(ctx context.Context, id string) (UserModel, error)
}

type UserRepo interface {
	Create(ctx context.Context, user UserModel) (UserModel, error)
	GetByUsername(ctx context.Context, username string) (UserModel, error)
	GetByEmail(ctx context.Context, email string) (UserModel, error)
	GetByID(ctx context.Context, id string) (UserModel, error)
	GetByIds(ctx context.Context, ids []string) ([]UserModel, error)
}

type UserModel struct {
	ID        string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
