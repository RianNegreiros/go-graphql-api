//go:build integration

package domain

import (
	"context"
	"github.com/RianNegreiros/go-graphql-api/internal/jwt"
	"log"
	"os"
	"testing"

	"github.com/RianNegreiros/go-graphql-api/config"
	"github.com/RianNegreiros/go-graphql-api/internal/domain"
	"github.com/RianNegreiros/go-graphql-api/internal/postgres"
)

var (
	conf             *config.Config
	db               *postgres.DB
	authService      *domain.AuthService
	userRepo         *postgres.UserRepo
	authTokenService *jwt.TokenService
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	config.LoadEnv(".env.test")

	conf = config.New()

	db = postgres.New(ctx, conf)
	defer db.Close()

	if err := db.Drop(); err != nil {
		log.Fatal(err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	userRepo = postgres.NewUserRepo(db)

	authTokenService = jwt.NewTokenService(conf)

	authService = domain.NewAuthService(userRepo, authTokenService)

	os.Exit(m.Run())
}
