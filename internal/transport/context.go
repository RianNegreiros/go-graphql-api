package transport

import (
	"context"

	"github.com/RianNegreiros/go-graphql-api/internal/user"
)

type contextKey string

var (
	ContextAuthIDKey contextKey = "currentUserId"
)

func GetUserIDFromContext(ctx context.Context) (string, error) {
	if ctx.Value(ContextAuthIDKey) == nil {
		return "", user.ErrNoUserIDInContext
	}

	userID, ok := ctx.Value(ContextAuthIDKey).(string)
	if !ok {
		return "", user.ErrNoUserIDInContext
	}

	return userID, nil
}

func PutUserIDIntoContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ContextAuthIDKey, id)
}
