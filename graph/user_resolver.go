package graph

import (
	"context"
	"fmt"
	"github.com/RianNegreiros/go-graphql-api/models"
)

func mapUser(user models.User) *User {
	return &User{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
}

func (r *queryResolver) Me(ctx context.Context) (*User, error) {
	panic(fmt.Errorf("not implemented"))
}
