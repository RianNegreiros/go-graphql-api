package test_helpers

import (
	"context"
	"testing"

	"github.com/RianNegreiros/go-graphql-api/internal/post"
	"github.com/RianNegreiros/go-graphql-api/internal/postgres"
	"github.com/RianNegreiros/go-graphql-api/internal/transport"
	"github.com/RianNegreiros/go-graphql-api/internal/user"
	"github.com/RianNegreiros/go-graphql-api/tests/faker"
	"github.com/stretchr/testify/require"
)

func TeardownDB(ctx context.Context, t *testing.T, db *postgres.DB) {
	t.Helper()

	err := db.Truncate(ctx)
	require.NoError(t, err)
}

func CreateUser(ctx context.Context, t *testing.T, userRepo user.UserRepo) user.UserModel {
	t.Helper()

	user, err := userRepo.Create(ctx, user.UserModel{
		Username: faker.Username(),
		Email:    faker.Email(),
		Password: faker.Password,
	})
	require.NoError(t, err)

	return user
}

func CreatePost(ctx context.Context, t *testing.T, postRepo post.PostRepo, forUser string) post.Post {
	t.Helper()

	post, err := postRepo.Create(ctx, post.Post{
		Body:   faker.RandStr(20),
		UserID: forUser,
	})
	require.NoError(t, err)

	return post
}

func LoginUser(ctx context.Context, t *testing.T, user user.UserModel) context.Context {
	t.Helper()

	return transport.PutUserIDIntoContext(ctx, user.ID)
}
