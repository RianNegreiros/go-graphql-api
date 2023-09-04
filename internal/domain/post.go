package domain

import (
	"context"

	"github.com/RianNegreiros/go-graphql-api/internal/post"
	"github.com/RianNegreiros/go-graphql-api/internal/transport"
	"github.com/RianNegreiros/go-graphql-api/internal/user"
	"github.com/RianNegreiros/go-graphql-api/internal/uuid"
)

type PostService struct {
	PostRepo post.PostRepo
}

func NewPostService(tr post.PostRepo) *PostService {
	return &PostService{
		PostRepo: tr,
	}
}

func (ts *PostService) All(ctx context.Context) ([]post.Post, error) {
	return ts.PostRepo.All(ctx)
}

func (ts *PostService) Create(ctx context.Context, input post.CreatePostInput) (post.Post, error) {
	currentUserID, err := transport.GetUserIDFromContext(ctx)
	if err != nil {
		return post.Post{}, user.ErrUnauthenticated
	}

	input.Sanitize()

	if err := input.Validate(); err != nil {
		return post.Post{}, err
	}

	p, err := ts.PostRepo.Create(ctx, post.Post{
		Body:   input.Body,
		UserID: currentUserID,
	})
	if err != nil {
		return post.Post{}, err
	}

	return p, nil
}

func (ts *PostService) GetByID(ctx context.Context, id string) (post.Post, error) {
	if !uuid.Validate(id) {
		return post.Post{}, uuid.ErrInvalidUUID
	}

	return ts.PostRepo.GetByID(ctx, id)
}

func (ts *PostService) Delete(ctx context.Context, id string) error {
	currentUserID, err := transport.GetUserIDFromContext(ctx)
	if err != nil {
		return user.ErrUnauthenticated
	}

	if !uuid.Validate(id) {
		return uuid.ErrInvalidUUID
	}

	post, err := ts.PostRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !post.CanDelete(user.UserModel{ID: currentUserID}) {
		return user.ErrForbidden
	}

	return ts.PostRepo.Delete(ctx, id)
}

func (ts *PostService) CreateReply(ctx context.Context, parentID string, input post.CreatePostInput) (post.Post, error) {
	currentUserID, err := transport.GetUserIDFromContext(ctx)
	if err != nil {
		return post.Post{}, user.ErrUnauthenticated
	}

	input.Sanitize()

	if err := input.Validate(); err != nil {
		return post.Post{}, err
	}

	if !uuid.Validate(parentID) {
		return post.Post{}, uuid.ErrInvalidUUID
	}

	if _, err := ts.PostRepo.GetByID(ctx, parentID); err != nil {
		return post.Post{}, user.ErrNotFound
	}

	p, err := ts.PostRepo.Create(ctx, post.Post{
		Body:     input.Body,
		UserID:   currentUserID,
		ParentID: &parentID,
	})
	if err != nil {
		return post.Post{}, err
	}

	return p, nil
}
