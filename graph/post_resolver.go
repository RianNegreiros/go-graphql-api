package graph

import (
	"context"

	"github.com/RianNegreiros/go-graphql-api/internal/post"
)

func mapPost(t post.Post) *Post {
	return &Post{
		ID:        t.ID,
		Body:      t.Body,
		UserID:    t.UserID,
		CreatedAt: t.CreatedAt,
	}
}

func mapPosts(posts []post.Post) []*Post {
	tt := make([]*Post, len(posts))

	for i, t := range posts {
		tt[i] = mapPost(t)
	}

	return tt
}

func (q *queryResolver) Posts(ctx context.Context) ([]*Post, error) {
	posts, err := q.PostService.All(ctx)
	if err != nil {
		return nil, err
	}

	return mapPosts(posts), nil
}

func (m *mutationResolver) CreatePost(ctx context.Context, input CreatePostInput) (*Post, error) {
	p, err := m.PostService.Create(ctx, post.CreatePostInput{
		Body: input.Body,
	})
	if err != nil {
		return nil, buildError(ctx, err)
	}

	return mapPost(p), nil
}

func (m *mutationResolver) DeletePost(ctx context.Context, id string) (bool, error) {
	if err := m.PostService.Delete(ctx, id); err != nil {
		return false, buildError(ctx, err)
	}

	return true, nil
}

func (t *postResolver) User(ctx context.Context, obj *Post) (*User, error) {
	user, err := t.UserService.GetByID(ctx, obj.UserID)
	if err != nil {
		return nil, buildError(ctx, err)
	}

	return mapUser(user), nil
}

func (m *mutationResolver) CreateReply(ctx context.Context, parentID string, input CreatePostInput) (*Post, error) {
	p, err := m.PostService.CreateReply(ctx, parentID, post.CreatePostInput{
		Body: input.Body,
	})
	if err != nil {
		return nil, buildError(ctx, err)
	}

	return mapPost(p), nil
}
