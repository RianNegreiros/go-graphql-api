package main

import (
	"net/http"

	"github.com/RianNegreiros/go-graphql-api/internal/transport"
	"github.com/RianNegreiros/go-graphql-api/internal/user"
)

func authMiddleware(service user.AuthTokenService) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token, err := service.ParseTokenFromRequest(ctx, r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx = transport.PutUserIDIntoContext(ctx, token.Sub)

			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
