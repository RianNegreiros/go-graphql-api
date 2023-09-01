package models

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrValidation         = errors.New("validation error")
	ErrNotFound           = errors.New("not found")
)

var (
	UsernameMinLength = 3
	PasswordMinLength = 6
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (AuthResponse, error)
	Login(ctx context.Context, input LoginInput) (AuthResponse, error)
}

type LoginInput struct {
	Email    string
	Password string
}

func (i *LoginInput) Sanitize() {
	i.Email = strings.TrimSpace(i.Email)
	i.Email = strings.ToLower(i.Email)
	i.Password = strings.TrimSpace(i.Password)
}

func (i LoginInput) Validate() error {
	if _, err := mail.ParseAddress(i.Email); err != nil {
		return fmt.Errorf("%w: invalid email address", ErrValidation)
	}

	if len(i.Password) < PasswordMinLength {
		return fmt.Errorf("%w: password required", ErrValidation)
	}

	return nil
}

type RegisterInput struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
}

type AuthResponse struct {
	AccessToken string
	User        User
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
