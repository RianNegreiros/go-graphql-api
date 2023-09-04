package domain

import (
	"context"

	"github.com/RianNegreiros/go-graphql-api/internal/user"
	"github.com/RianNegreiros/go-graphql-api/internal/uuid"
)

type UserService struct {
	UserRepo user.UserRepo
}

func NewUserService(ur user.UserRepo) *UserService {
	return &UserService{
		UserRepo: ur,
	}
}

func (u *UserService) GetByID(ctx context.Context, id string) (user.UserModel, error) {
	if !uuid.Validate(id) {
		return user.UserModel{}, uuid.ErrInvalidUUID
	}

	return u.UserRepo.GetByID(ctx, id)
}
