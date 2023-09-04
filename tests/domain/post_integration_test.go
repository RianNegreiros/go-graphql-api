//go:build integration
// +build integration

package domain

import (
	"context"
	"testing"

	"github.com/RianNegreiros/go-graphql-api/internal/post"
	"github.com/RianNegreiros/go-graphql-api/internal/user"
	"github.com/RianNegreiros/go-graphql-api/internal/uuid"
	"github.com/RianNegreiros/go-graphql-api/tests/faker"
	"github.com/RianNegreiros/go-graphql-api/tests/test_helpers"
	"github.com/stretchr/testify/require"
)

func TestIntegrationPostService_Create(t *testing.T) {
	t.Run("not auth user cannot create a post", func(t *testing.T) {
		ctx := context.Background()

		_, err := postService.Create(ctx, post.CreatePostInput{
			Body: "test",
		})

		require.ErrorIs(t, err, user.ErrUnauthenticated)
	})

	t.Run("can create a post", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		currentUser := test_helpers.CreateUser(ctx, t, userRepo)

		ctx = test_helpers.LoginUser(ctx, t, currentUser)

		input := post.CreatePostInput{
			Body: faker.RandStr(100),
		}

		post, err := postService.Create(ctx, input)
		require.NoError(t, err)

		require.NotEmpty(t, post.ID, "post.ID")
		require.Equal(t, input.Body, post.Body, "post.Body")
		require.Equal(t, currentUser.ID, post.UserID, "post.UserID")
		require.NotEmpty(t, post.CreatedAt, "post.CreatedAt")
	})
}

func TestIntegrationPostService_All(t *testing.T) {
	t.Run("return all posts", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		user := test_helpers.CreateUser(ctx, t, userRepo)

		test_helpers.CreatePost(ctx, t, postRepo, user.ID)
		test_helpers.CreatePost(ctx, t, postRepo, user.ID)
		test_helpers.CreatePost(ctx, t, postRepo, user.ID)

		posts, err := postService.All(ctx)
		require.NoError(t, err)

		require.Len(t, posts, 3)
	})
}

func TestIntegrationPostService_GetByID(t *testing.T) {
	t.Run("can get a post by id", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		user := test_helpers.CreateUser(ctx, t, userRepo)
		existingPost := test_helpers.CreatePost(ctx, t, postRepo, user.ID)

		post, err := postService.GetByID(ctx, existingPost.ID)
		require.NoError(t, err)

		require.Equal(t, existingPost.ID, post.ID, "post.ID")
		require.Equal(t, existingPost.Body, post.Body, "post.Body")
	})

	t.Run("return error not found if the post doesn't exist", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		_, err := postService.GetByID(ctx, faker.UUID())
		require.ErrorIs(t, err, user.ErrNotFound)
	})

	t.Run("return error invalid uuid", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		_, err := postService.GetByID(ctx, "123")
		require.ErrorIs(t, err, uuid.ErrInvalidUUID)
	})
}

func TestIntegrationPostService_Delete(t *testing.T) {
	t.Run("not auth user cannot delete a post", func(t *testing.T) {
		ctx := context.Background()

		err := postService.Delete(ctx, faker.UUID())
		require.ErrorIs(t, err, user.ErrUnauthenticated)
	})

	t.Run("cannot delete a post if not the owner", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		otherUser := test_helpers.CreateUser(ctx, t, userRepo)
		currentUser := test_helpers.CreateUser(ctx, t, userRepo)

		post := test_helpers.CreatePost(ctx, t, postRepo, otherUser.ID)

		ctx = test_helpers.LoginUser(ctx, t, currentUser)

		_, err := postRepo.GetByID(ctx, post.ID)
		require.NoError(t, err)

		err = postService.Delete(ctx, post.ID)
		require.ErrorIs(t, err, user.ErrForbidden)

		_, err = postRepo.GetByID(ctx, post.ID)
		require.NoError(t, err)
	})

	t.Run("can delete a post", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		currentUser := test_helpers.CreateUser(ctx, t, userRepo)

		post := test_helpers.CreatePost(ctx, t, postRepo, currentUser.ID)

		ctx = test_helpers.LoginUser(ctx, t, currentUser)

		_, err := postRepo.GetByID(ctx, post.ID)
		require.NoError(t, err)

		err = postService.Delete(ctx, post.ID)
		require.NoError(t, err)

		_, err = postRepo.GetByID(ctx, post.ID)
		require.ErrorIs(t, err, user.ErrNotFound)
	})
}

func TestIntegrationPostService_CreateReply(t *testing.T) {
	t.Run("not auth user cannot create a reply to a post", func(t *testing.T) {
		ctx := context.Background()

		_, err := postService.CreateReply(ctx, faker.UUID(), post.CreatePostInput{
			Body: faker.RandStr(20),
		})
		require.ErrorIs(t, err, user.ErrUnauthenticated)
	})

	t.Run("cannot create a reply to a not found post", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		currentUser := test_helpers.CreateUser(ctx, t, userRepo)

		ctx = test_helpers.LoginUser(ctx, t, currentUser)

		_, err := postService.CreateReply(ctx, faker.UUID(), post.CreatePostInput{
			Body: faker.RandStr(20),
		})
		require.ErrorIs(t, err, user.ErrNotFound)
	})

	t.Run("can create a reply to a post", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		currentUser := test_helpers.CreateUser(ctx, t, userRepo)

		ctx = test_helpers.LoginUser(ctx, t, currentUser)

		p := test_helpers.CreatePost(ctx, t, postRepo, currentUser.ID)

		input := post.CreatePostInput{
			Body: faker.RandStr(20),
		}

		reply, err := postService.CreateReply(ctx, p.ID, input)
		require.NoError(t, err)

		require.NotEmpty(t, reply.ID, "reply.ID")
		require.Equal(t, input.Body, reply.Body, "reply.Body")
		require.Equal(t, currentUser.ID, reply.UserID, "reply.UserID")
		require.Equal(t, p.ID, *reply.ParentID, "reply.ParentID")
		require.NotEmpty(t, reply.CreatedAt, "reply.CreatedAt")
	})
}
