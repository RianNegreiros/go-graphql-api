package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/RianNegreiros/go-graphql-api/internal"
	"net/mail"
	"strings"
)

var (
	ErrValidation = errors.New("validation error")
)

var (
	UsernameMinLength = 3
	PasswordMinLength = 6
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (AuthResponse, error)
}

type RegisterInput struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
}

type AuthResponse struct {
	AccessToken string
	User        internal.User
}

func (r RegisterInput) Validate() error {
	if len(r.Username) < UsernameMinLength {
		return fmt.Errorf("%w: username must be at least %d characters long", ErrValidation, UsernameMinLength)
	}

	if _, err := mail.ParseAddress(r.Email); err != nil {
		return fmt.Errorf("%w: invalid email address", ErrValidation)
	}

	if len(r.Password) < PasswordMinLength {
		return fmt.Errorf("%w: password must be at least %d characters long", ErrValidation, PasswordMinLength)
	}

	if r.Password != r.ConfirmPassword {
		return fmt.Errorf("%w: password and confirm password must match", ErrValidation)
	}

	return nil
}

func (r *RegisterInput) Sanitize() {
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
	r.Password = strings.TrimSpace(r.Password)
	r.ConfirmPassword = strings.TrimSpace(r.ConfirmPassword)
}
