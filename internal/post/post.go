package post

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/RianNegreiros/go-graphql-api/internal/user"
)

var (
	PostMinLength = 2
	PostMaxLength = 250
)

type CreatePostInput struct {
	Body string
}

func (in *CreatePostInput) Sanitize() {
	in.Body = strings.TrimSpace(in.Body)
}

func (in CreatePostInput) Validate() error {
	if len(in.Body) < PostMinLength {
		return fmt.Errorf("%w: body not long enough, (%d) characters at least", user.ErrValidation, PostMinLength)
	}

	if len(in.Body) > PostMaxLength {
		return fmt.Errorf("%w: body too long, (%d) characters at max", user.ErrValidation, PostMaxLength)
	}

	return nil
}

type Post struct {
	ID        string
	Body      string
	UserID    string
	ParentID  *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t Post) CanDelete(user user.UserModel) bool {
	return t.UserID == user.ID
}

type PostService interface {
	All(ctx context.Context) ([]Post, error)
	Create(ctx context.Context, input CreatePostInput) (Post, error)
	CreateReply(ctx context.Context, parentID string, input CreatePostInput) (Post, error)
	GetByID(ctx context.Context, id string) (Post, error)
	Delete(ctx context.Context, id string) error
}

type PostRepo interface {
	All(ctx context.Context) ([]Post, error)
	Create(ctx context.Context, Post Post) (Post, error)
	GetByID(ctx context.Context, id string) (Post, error)
	Delete(ctx context.Context, id string) error
}
