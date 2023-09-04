package transport

import (
	"context"
	"testing"

	"github.com/RianNegreiros/go-graphql-api/internal/transport"
	"github.com/RianNegreiros/go-graphql-api/internal/user"
	"github.com/stretchr/testify/require"
)

func TestGetUserIDFromContext(t *testing.T) {
	t.Run("should return user id from context", func(t *testing.T) {
		ctx := context.Background()

		ctx = context.WithValue(ctx, transport.ContextAuthIDKey, "1")

		userID, err := transport.GetUserIDFromContext(ctx)
		require.NoError(t, err)
		require.Equal(t, "1", userID)
	})

	t.Run("return error if no id", func(t *testing.T) {
		ctx := context.Background()

		_, err := transport.GetUserIDFromContext(ctx)
		require.ErrorIs(t, err, user.ErrNoUserIDInContext)
	})

	t.Run("return error if value is not a string", func(t *testing.T) {
		ctx := context.Background()

		ctx = context.WithValue(ctx, transport.ContextAuthIDKey, 123)

		_, err := transport.GetUserIDFromContext(ctx)
		require.ErrorIs(t, err, user.ErrNoUserIDInContext)

	})
}

func TestPutUserIDIntoContext(t *testing.T) {
	t.Run("add user id into context", func(t *testing.T) {
		ctx := context.Background()

		ctx = transport.PutUserIDIntoContext(ctx, "123")

		userID, err := transport.GetUserIDFromContext(ctx)
		require.NoError(t, err)
		require.Equal(t, "123", userID)
	})
}
