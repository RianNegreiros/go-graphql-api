package main

import (
	"github.com/RianNegreiros/go-graphql-api/models"
	"github.com/RianNegreiros/go-graphql-api/transport"
	"net/http"
)

func authMiddleware(service models.AuthTokenService) func(handler http.Handler) http.Handler {
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
