package graph

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/RianNegreiros/go-graphql-api/models"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"net/http"
)

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	AuthService models.AuthService
}

type queryResolver struct{ *Resolver }

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func buildBadRequestError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": http.StatusBadRequest,
		},
	}
}
