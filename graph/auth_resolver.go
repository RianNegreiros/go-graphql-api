package graph

import (
	"context"
	"errors"
	"github.com/RianNegreiros/go-graphql-api/models"
)

func mapAuthResponse(authResponse models.AuthResponse) *AuthResponse {
	return &AuthResponse{
		AccessToken: authResponse.AccessToken,
		User:        mapUser(authResponse.User),
	}
}

func (m *mutationResolver) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	res, err := m.AuthService.Register(ctx, models.RegisterInput{
		Email:           input.Email,
		Username:        input.Username,
		Password:        input.Password,
		ConfirmPassword: input.ConfirmPassword,
	})
	if err != nil {
		switch {
		case errors.Is(err, models.ErrValidation) ||
			errors.Is(err, models.ErrEmailTaken) ||
			errors.Is(err, models.ErrUsernameTaken):
			return nil, buildBadRequestError(ctx, err)
		default:
			return nil, err
		}
	}

	return mapAuthResponse(res), nil
}

func (m *mutationResolver) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	res, err := m.AuthService.Login(ctx, models.LoginInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, models.ErrValidation) ||
			errors.Is(err, models.ErrInvalidCredentials):
			return nil, buildBadRequestError(ctx, err)
		default:
			return nil, err
		}
	}

	return mapAuthResponse(res), nil
}
