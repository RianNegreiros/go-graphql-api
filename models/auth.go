package models

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrValidation         = errors.New("validation error")
	ErrNotFound           = errors.New("not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrNoUserIDInContext  = errors.New("no user id in context")
)

var (
	UsernameMinLength = 3
	PasswordMinLength = 6
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (AuthResponse, error)
	Login(ctx context.Context, input LoginInput) (AuthResponse, error)
}
type AuthTokenService interface {
	CreateAccessToken(ctx context.Context, user User) (string, error)
	CreateRefreshToken(ctx context.Context, user User, tokenID string) (string, error)
	ParseToken(ctx context.Context, payload string) (AuthToken, error)
	ParseTokenFromRequest(ctx context.Context, r *http.Request) (AuthToken, error)
}

type AuthResponse struct {
	AccessToken string
	User        User
}

type RegisterInput struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthToken struct {
	ID  string
	Sub string
}

func (i LoginInput) Sanitize() LoginInput {
	i.Email = strings.TrimSpace(i.Email)
	i.Email = strings.ToLower(i.Email)

	i.Password = strings.TrimSpace(i.Password)

	return i
}

func (i LoginInput) Validate() error {
	if _, err := mail.ParseAddress(i.Email); err != nil {
		return fmt.Errorf("%w: invalid email address", ErrValidation)
	}

	if len(i.Password) < 1 {
		return fmt.Errorf("%w: password required", ErrValidation)
	}

	return nil
}

func (i RegisterInput) Validate() error {
	if len(i.Username) < UsernameMinLength {
		return fmt.Errorf("%w: username must be at least %d characters long", ErrValidation, UsernameMinLength)
	}

	if _, err := mail.ParseAddress(i.Email); err != nil {
		return fmt.Errorf("%w: invalid email address", ErrValidation)
	}

	if len(i.Password) < PasswordMinLength {
		return fmt.Errorf("%w: password must be at least %d characters long", ErrValidation, PasswordMinLength)
	}

	if i.Password != i.ConfirmPassword {
		return fmt.Errorf("%w: password and confirm password must match", ErrValidation)
	}

	return nil
}

func (i RegisterInput) Sanitize() RegisterInput {
	i.Email = strings.TrimSpace(i.Email)
	i.Email = strings.ToLower(i.Email)

	i.Username = strings.TrimSpace(i.Username)
	i.Password = strings.TrimSpace(i.Password)
	i.ConfirmPassword = strings.TrimSpace(i.ConfirmPassword)

	return i
}
