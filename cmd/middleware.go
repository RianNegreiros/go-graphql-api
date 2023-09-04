package main

import (
	"github.com/RianNegreiros/go-graphql-api/internal/transport"
	"github.com/RianNegreiros/go-graphql-api/internal/user"
	"net/http"
)

func authMiddleware(authTokenService user.AuthTokenService) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token, err := authTokenService.ParseTokenFromRequest(ctx, r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx = transport.PutUserIDIntoContext(ctx, token.Sub)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
