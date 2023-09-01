package tests

import (
	"context"
	"github.com/RianNegreiros/go-graphql-api/internal/postgres"
	"github.com/stretchr/testify/require"
	"testing"
)

func TeardownDB(ctx context.Context, t *testing.T, db *postgres.DB) {
	t.Helper()

	err := db.Truncate(ctx)
	require.NoError(t, err)
}
