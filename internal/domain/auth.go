package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/RianNegreiros/go-graphql-api/internal/user"
	"golang.org/x/crypto/bcrypt"
)

var passwordCost = bcrypt.DefaultCost

type AuthService struct {
	AuthTokenService user.AuthTokenService
	UserRepo         user.UserRepo
}

func NewAuthService(ur user.UserRepo, service user.AuthTokenService) *AuthService {
	return &AuthService{
		UserRepo:         ur,
		AuthTokenService: service,
	}
}

func (as *AuthService) Register(ctx context.Context, input user.RegisterInput) (user.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return user.AuthResponse{}, err
	}

	if _, err := as.UserRepo.GetByUsername(ctx, input.Username); !errors.Is(err, user.ErrNotFound) {
		return user.AuthResponse{}, user.ErrUsernameTaken
	}

	if _, err := as.UserRepo.GetByEmail(ctx, input.Email); !errors.Is(err, user.ErrNotFound) {
		return user.AuthResponse{}, user.ErrEmailTaken
	}

	u := user.UserModel{
		Email:    input.Email,
		Username: input.Username,
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), passwordCost)
	if err != nil {
		return user.AuthResponse{}, fmt.Errorf("error hashing password: %v", err)
	}

	u.Password = string(hashPassword)

	u, err = as.UserRepo.Create(ctx, u)
	if err != nil {
		return user.AuthResponse{}, fmt.Errorf("error creating user: %v", err)
	}

	accessToken, err := as.AuthTokenService.CreateAccessToken(ctx, u)
	if err != nil {
		return user.AuthResponse{}, user.ErrGenerateToken
	}

	return user.AuthResponse{
		AccessToken: accessToken,
		User:        u,
	}, nil
}

func (as *AuthService) Login(ctx context.Context, input user.LoginInput) (user.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return user.AuthResponse{}, err
	}

	u, err := as.UserRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return user.AuthResponse{}, user.ErrInvalidCredentials
		default:
			return user.AuthResponse{}, err
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		return user.AuthResponse{}, user.ErrInvalidCredentials
	}

	accessToken, err := as.AuthTokenService.CreateAccessToken(ctx, u)
	if err != nil {
		return user.AuthResponse{}, user.ErrGenerateToken
	}

	return user.AuthResponse{
		AccessToken: accessToken,
		User:        u,
	}, nil
}
