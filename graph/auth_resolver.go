package graph

import (
	"context"
	"fmt"
)

func (r *mutationResolver) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	panic(fmt.Errorf("not implemented"))
}
