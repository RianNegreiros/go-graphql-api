package graph

import (
	"context"
	"errors"
	"fmt"

	"github.com/RianNegreiros/go-graphql-api/internal/user"
)

func mapAuthResponse(a user.AuthResponse) *AuthResponse {
	return &AuthResponse{
		AccessToken: a.AccessToken,
		User:        mapUser(a.User),
	}
}

func (m *mutationResolver) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	res, err := m.AuthService.Register(ctx, user.RegisterInput{
		Email:           input.Email,
		Username:        input.Username,
		Password:        input.Password,
		ConfirmPassword: input.ConfirmPassword,
	})

	fmt.Println(res)
	fmt.Println(input)

	if err != nil {
		switch {
		case errors.Is(err, user.ErrValidation) ||
			errors.Is(err, user.ErrEmailTaken) ||
			errors.Is(err, user.ErrUsernameTaken):
			return nil, buildBadRequestError(ctx, err)
		default:
			return nil, err
		}
	}

	return mapAuthResponse(res), nil
}

func (m *mutationResolver) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	res, err := m.AuthService.Login(ctx, user.LoginInput{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		switch {
		case errors.Is(err, user.ErrValidation) ||
			errors.Is(err, user.ErrInvalidCredentials):
			return nil, buildBadRequestError(ctx, err)
		default:
			return nil, err
		}
	}

	return mapAuthResponse(res), nil
}
