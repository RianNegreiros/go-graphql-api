package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/RianNegreiros/go-graphql-api/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo models.UserRepo
}

func NewAuthService(userRepo models.UserRepo) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, input models.RegisterInput) (models.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return models.AuthResponse{}, err
	}

	if _, err := s.UserRepo.GetByUsername(ctx, input.Username); !errors.Is(err, models.ErrNotFound) {
		return models.AuthResponse{}, models.ErrUsernameTaken
	}

	if _, err := s.UserRepo.GetByEmail(ctx, input.Email); !errors.Is(err, models.ErrNotFound) {
		return models.AuthResponse{}, models.ErrEmailTaken
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return models.AuthResponse{}, fmt.Errorf("%w: error generating password hash", err)
	}

	user.Password = string(hashPassword)

	user, err = s.UserRepo.Create(ctx, user)
	if err != nil {
		return models.AuthResponse{}, fmt.Errorf("%w: error creating user", err)
	}

	return models.AuthResponse{
		AccessToken: "access_token",
		User:        user,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, input models.LoginInput) (models.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return models.AuthResponse{}, err
	}

	user, err := s.UserRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return models.AuthResponse{}, models.ErrInvalidCredentials
		default:
			return models.AuthResponse{}, err
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return models.AuthResponse{}, models.ErrInvalidCredentials
	}

	return models.AuthResponse{
		AccessToken: "access_token",
		User:        user,
	}, nil
}
