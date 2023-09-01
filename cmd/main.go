package main

import (
	"context"
	"log"

	"github.com/RianNegreiros/go-graphql-api/config"
	"github.com/RianNegreiros/go-graphql-api/postgres"
)

func main() {
	ctx := context.Background()

	config := config.New()

	db := postgres.New(ctx, config)

	if err := db.Migrate(); err != nil {
		log.Fatalf("error migrating postgres: %v", err)
	}
}
