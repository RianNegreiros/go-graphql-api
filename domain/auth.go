package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/RianNegreiros/go-graphql-api/models"
	"golang.org/x/crypto/bcrypt"
)

var passwordCost = bcrypt.DefaultCost

type AuthService struct {
	UserRepo models.UserRepo
}

func NewAuthService(ur models.UserRepo) *AuthService {
	return &AuthService{
		UserRepo: ur,
	}
}

func (as *AuthService) Register(ctx context.Context, input models.RegisterInput) (models.AuthResponse, error) {
	input = input.Sanitize()

	if err := input.Validate(); err != nil {
		return models.AuthResponse{}, err
	}

	if _, err := as.UserRepo.GetByUsername(ctx, input.Username); !errors.Is(err, models.ErrNotFound) {
		return models.AuthResponse{}, models.ErrUsernameTaken
	}

	if _, err := as.UserRepo.GetByEmail(ctx, input.Email); !errors.Is(err, models.ErrNotFound) {
		return models.AuthResponse{}, models.ErrEmailTaken
	}

	user := models.User{
		Email:    input.Email,
		Username: input.Username,
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), passwordCost)
	if err != nil {
		return models.AuthResponse{}, fmt.Errorf("error hashing password: %v", err)
	}

	user.Password = string(hashPassword)

	user, err = as.UserRepo.Create(ctx, user)
	if err != nil {
		return models.AuthResponse{}, fmt.Errorf("error creating user: %v", err)
	}

	return models.AuthResponse{
		AccessToken: "access_token",
		User:        user,
	}, nil
}

func (as *AuthService) Login(ctx context.Context, input models.LoginInput) (models.AuthResponse, error) {
	input = input.Sanitize()

	if err := input.Validate(); err != nil {
		return models.AuthResponse{}, err
	}

	user, err := as.UserRepo.GetByEmail(ctx, input.Email)
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
