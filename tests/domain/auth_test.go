package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/RianNegreiros/go-graphql-api/internal/domain"
	"github.com/RianNegreiros/go-graphql-api/internal/user"
	mocks "github.com/RianNegreiros/go-graphql-api/mocks/internal_/user"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	validInput := user.RegisterInput{
		Username:        "john",
		Email:           "johndoe@mail.com",
		Password:        "123456",
		ConfirmPassword: "123456",
	}

	t.Run("valid input", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByUsername", mock.Anything, mock.Anything).
			Return(user.UserModel{}, user.ErrNotFound)

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).
			Return(user.UserModel{}, user.ErrNotFound)

		userRepo.On("Create", mock.Anything, mock.Anything).
			Return(user.UserModel{
				ID:       "user_id",
				Username: validInput.Username,
				Email:    validInput.Email,
				Password: "hashed_password",
			}, nil)

		service := domain.NewAuthService(userRepo)

		res, err := service.Register(ctx, validInput)
		require.NoError(t, err)

		require.NotEmpty(t, res.AccessToken)
		require.NotEmpty(t, res.User.ID)
		require.Equal(t, validInput.Username, res.User.Username)
		require.Equal(t, validInput.Email, res.User.Email)
		require.NotEmpty(t, res.User.Password)

		userRepo.AssertExpectations(t)
	})

	t.Run("username taken", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByUsername", mock.Anything, mock.Anything).
			Return(user.UserModel{}, nil)

		service := domain.NewAuthService(userRepo)

		_, err := service.Register(ctx, validInput)
		require.ErrorIs(t, err, user.ErrUsernameTaken)

		userRepo.AssertNotCalled(t, "Create")
		userRepo.AssertExpectations(t)
	})

	t.Run("email taken", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByUsername", mock.Anything, mock.Anything).
			Return(user.UserModel{}, user.ErrNotFound)

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).
			Return(user.UserModel{}, nil)

		service := domain.NewAuthService(userRepo)

		_, err := service.Register(ctx, validInput)
		require.ErrorIs(t, err, user.ErrEmailTaken)

		userRepo.AssertNotCalled(t, "Create")
		userRepo.AssertExpectations(t)
	})

	t.Run("error creating user", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByUsername", mock.Anything, mock.Anything).
			Return(user.UserModel{}, user.ErrNotFound)

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).
			Return(user.UserModel{}, user.ErrNotFound)

		userRepo.On("Create", mock.Anything, mock.Anything).
			Return(user.UserModel{}, errors.New("some error"))

		service := domain.NewAuthService(userRepo)

		_, err := service.Register(ctx, validInput)
		require.Error(t, err)

		userRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}

		service := domain.NewAuthService(userRepo)

		_, err := service.Register(ctx, user.RegisterInput{})
		require.Error(t, err)

		userRepo.AssertNotCalled(t, "GetByUsername")
		userRepo.AssertNotCalled(t, "GetByEmail")
		userRepo.AssertNotCalled(t, "Create")
		userRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	validInput := user.LoginInput{
		Email:    "johndoe@mail.com",
		Password: "hashed_password",
	}

	t.Run("valid input", func(t *testing.T) {
		ctx := context.Background()

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(validInput.Password), bcrypt.DefaultCost)
		require.NoError(t, err)

		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).
			Return(user.UserModel{
				ID:       "user_id",
				Username: "john",
				Email:    validInput.Email,
				Password: string(hashedPassword),
			}, nil)

		service := domain.NewAuthService(userRepo)

		res, err := service.Login(ctx, validInput)
		require.NoError(t, err)

		require.NotEmpty(t, res.AccessToken)
		require.NotEmpty(t, res.User.ID)
		require.NotEmpty(t, res.User.Username)
		require.Equal(t, validInput.Email, res.User.Email)

		userRepo.AssertExpectations(t)
	})

	t.Run("invalid email", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).
			Return(user.UserModel{}, user.ErrNotFound)

		service := domain.NewAuthService(userRepo)

		_, err := service.Login(ctx, validInput)
		require.ErrorIs(t, err, user.ErrInvalidCredentials)

		userRepo.AssertExpectations(t)
	})

	t.Run("get user by email error", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).
			Return(user.UserModel{}, errors.New("some error"))

		service := domain.NewAuthService(userRepo)

		_, err := service.Login(ctx, validInput)
		require.Error(t, err)

		userRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}

		service := domain.NewAuthService(userRepo)

		_, err := service.Login(ctx, user.LoginInput{})
		require.ErrorIs(t, err, user.ErrValidation)

		userRepo.AssertNotCalled(t, "GetByEmail")
	})
}
