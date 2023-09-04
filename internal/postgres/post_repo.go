package postgres

import (
	"context"
	"fmt"

	"github.com/RianNegreiros/go-graphql-api/internal/post"
	"github.com/RianNegreiros/go-graphql-api/internal/user"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type PostRepo struct {
	DB *DB
}

func NewPostRepo(db *DB) *PostRepo {
	return &PostRepo{
		DB: db,
	}
}

func (tr *PostRepo) All(ctx context.Context) ([]post.Post, error) {
	return getAllPost(ctx, tr.DB.Pool)
}

func getAllPost(ctx context.Context, q pgxscan.Querier) ([]post.Post, error) {
	query := `SELECT * FROM posts ORDER BY created_at DESC;`

	var posts []post.Post

	if err := pgxscan.Select(ctx, q, &posts, query); err != nil {
		return nil, fmt.Errorf("error get all posts %+v", err)
	}

	return posts, nil
}

func (tr *PostRepo) Create(ctx context.Context, p post.Post) (post.Post, error) {
	tx, err := tr.DB.Pool.Begin(ctx)
	if err != nil {
		return post.Post{}, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	p, err = createPost(ctx, tx, p)
	if err != nil {
		return post.Post{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return post.Post{}, fmt.Errorf("error commiting: %v", err)
	}

	return p, nil
}

func createPost(ctx context.Context, tx pgx.Tx, p post.Post) (post.Post, error) {
	query := `INSERT INTO posts (body, user_id, parent_id) VALUES ($1, $2, $3) RETURNING *;`

	t := post.Post{}

	if err := pgxscan.Get(ctx, tx, &t, query, p.Body, p.UserID, p.ParentID); err != nil {
		return post.Post{}, fmt.Errorf("error insert: %v", err)
	}

	return t, nil
}

func (tr *PostRepo) GetByID(ctx context.Context, id string) (post.Post, error) {
	return getPostByID(ctx, tr.DB.Pool, id)
}

func getPostByID(ctx context.Context, q pgxscan.Querier, id string) (post.Post, error) {
	query := `SELECT * FROM posts WHERE id = $1 LIMIT 1;`

	t := post.Post{}

	if err := pgxscan.Get(ctx, q, &t, query, id); err != nil {
		if pgxscan.NotFound(err) {
			return post.Post{}, user.ErrNotFound
		}

		return post.Post{}, fmt.Errorf("error get post: %+v", err)
	}

	return t, nil
}

func (tr *PostRepo) Delete(ctx context.Context, id string) error {
	tx, err := tr.DB.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	if err := deletePost(ctx, tx, id); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error commiting: %v", err)
	}

	return nil
}

func deletePost(ctx context.Context, tx pgx.Tx, id string) error {
	query := `DELETE FROM posts WHERE id = $1;`

	if _, err := tx.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("error insert: %v", err)
	}

	return nil
}
