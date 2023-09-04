package graph

import (
	"context"
	"github.com/RianNegreiros/go-graphql-api/internal/transport"
	"github.com/RianNegreiros/go-graphql-api/internal/user"
)

func mapUser(user user.UserModel) *User {
	return &User{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
}

func (r *queryResolver) Me(ctx context.Context) (*User, error) {
	userID, err := transport.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, user.ErrUnauthenticated
	}

	return mapUser(user.UserModel{
		ID: userID,
	}), nil
}
