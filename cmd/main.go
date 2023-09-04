package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/RianNegreiros/go-graphql-api/config"
	"github.com/RianNegreiros/go-graphql-api/graph"
	"github.com/RianNegreiros/go-graphql-api/internal/domain"
	"github.com/RianNegreiros/go-graphql-api/internal/jwt"
	"github.com/RianNegreiros/go-graphql-api/internal/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx := context.Background()

	config.LoadEnv(".env")

	conf := config.New()

	db := postgres.New(ctx, conf)

	if err := db.Migrate(); err != nil {
		log.Fatalf("error migrating postgres: %v", err)
	}

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RedirectSlashes)
	router.Use(middleware.Timeout(time.Second * 60))

	userRepo := postgres.NewUserRepo(db)

	authTokenService := jwt.NewTokenService(conf)
	authService := domain.NewAuthService(userRepo)

	router.Use(authMiddleware(authTokenService))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					AuthService: authService,
				},
			},
		),
	))

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
