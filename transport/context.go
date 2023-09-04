package transport

import (
	"context"
	"github.com/RianNegreiros/go-graphql-api/models"
)

type contextKey string

var (
	ContextAuthIDKey contextKey = "currentUserId"
)

func GetUserIDFromContext(ctx context.Context) (string, error) {
	if ctx.Value(ContextAuthIDKey) == nil {
		return "", models.ErrNoUserIDInContext
	}

	userID, ok := ctx.Value(ContextAuthIDKey).(string)
	if !ok {
		return "", models.ErrNoUserIDInContext
	}

	return userID, nil
}
