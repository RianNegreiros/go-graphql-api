//go:build integration

package domain

import (
	"context"
	"testing"

	"github.com/RianNegreiros/go-graphql-api/internal/user"
	"github.com/RianNegreiros/go-graphql-api/tests/test_helpers"
	"github.com/stretchr/testify/require"
)

func TestIntegrationAuthService_Register(t *testing.T) {
	validInput := user.RegisterInput{
		Username:        "john",
		Email:           "johndoe@mail.com",
		Password:        "123456",
		ConfirmPassword: "123456",
	}

	t.Run("valid input", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		res, err := authService.Register(ctx, validInput)
		require.NoError(t, err)

		require.NotEmpty(t, res.AccessToken)
		require.NotEmpty(t, res.User.ID)
		require.Equal(t, validInput.Username, res.User.Username)
		require.Equal(t, validInput.Email, res.User.Email)
		require.NotEmpty(t, res.User.Password)
		require.NotEqual(t, validInput.Password, res.User.Password)
	})

	t.Run("username taken", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		_, err := authService.Register(ctx, validInput)
		require.NoError(t, err)

		_, err = authService.Register(ctx, user.RegisterInput{
			Username:        validInput.Username,
			Email:           "johndoe@mail.com",
			Password:        "123456",
			ConfirmPassword: "123456",
		})
		require.ErrorIs(t, err, user.ErrUsernameTaken)
	})

	t.Run("email taken", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		_, err := authService.Register(ctx, validInput)
		require.NoError(t, err)

		_, err = authService.Register(ctx, user.RegisterInput{
			Username:        "john2",
			Email:           validInput.Email,
			Password:        "123456",
			ConfirmPassword: "123456",
		})

		require.ErrorIs(t, err, user.ErrEmailTaken)
	})
}
