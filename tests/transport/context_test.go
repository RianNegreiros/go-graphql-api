package transport

import (
	"context"
	"github.com/RianNegreiros/go-graphql-api/models"
	"github.com/RianNegreiros/go-graphql-api/transport"
	"github.com/stretchr/testify/require"
	"testing"
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
		require.ErrorIs(t, err, models.ErrNoUserIDInContext)
	})

	t.Run("return error if value is not a string", func(t *testing.T) {
		ctx := context.Background()

		ctx = context.WithValue(ctx, transport.ContextAuthIDKey, 123)

		_, err := transport.GetUserIDFromContext(ctx)
		require.ErrorIs(t, err, models.ErrNoUserIDInContext)

	})
}
