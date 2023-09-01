package tests

import (
	"context"
	"testing"

	"github.com/RianNegreiros/go-graphql-api/postgres"
	"github.com/stretchr/testify/require"
)

func TeardownDB(ctx context.Context, t *testing.T, db *postgres.DB) {
	t.Helper()

	err := db.Truncate(ctx)
	require.NoError(t, err)
}
