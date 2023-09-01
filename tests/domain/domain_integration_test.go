//go:build integration

package domain

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/RianNegreiros/go-graphql-api/config"
	"github.com/RianNegreiros/go-graphql-api/internal/domain"
	"github.com/RianNegreiros/go-graphql-api/internal/postgres"
)

var (
	conf        *config.Config
	db          *postgres.DB
	authService *domain.AuthService
	userRepo    *postgres.UserRepo
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

	authService = domain.NewAuthService(userRepo)

	os.Exit(m.Run())
}
