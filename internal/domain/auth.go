package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/RianNegreiros/go-graphql-api/internal"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo internal.UserRepo
}

func NewAuthService(userRepo internal.UserRepo) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, input internal.RegisterInput) (internal.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return internal.AuthResponse{}, err
	}

	if _, err := s.UserRepo.GetByUsername(ctx, input.Username); !errors.Is(err, internal.ErrNotFound) {
		return internal.AuthResponse{}, internal.ErrUsernameTaken
	}

	if _, err := s.UserRepo.GetByEmail(ctx, input.Email); !errors.Is(err, internal.ErrNotFound) {
		return internal.AuthResponse{}, internal.ErrEmailTaken
	}

	user := internal.User{
		Username: input.Username,
		Email:    input.Email,
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return internal.AuthResponse{}, fmt.Errorf("%w: error generating password hash", err)
	}

	user.Password = string(hashPassword)

	user, err = s.UserRepo.Create(ctx, user)
	if err != nil {
		return internal.AuthResponse{}, fmt.Errorf("%w: error creating user", err)
	}

	return internal.AuthResponse{
		AccessToken: "access_token",
		User:        user,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, input internal.LoginInput) (internal.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return internal.AuthResponse{}, err
	}

	user, err := s.UserRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrNotFound):
			return internal.AuthResponse{}, internal.ErrInvalidCredentials
		default:
			return internal.AuthResponse{}, err
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return internal.AuthResponse{}, internal.ErrInvalidCredentials
	}

	return internal.AuthResponse{
		AccessToken: "access_token",
		User:        user,
	}, nil
}
