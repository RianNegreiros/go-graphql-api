package main

import (
	"context"
	"github.com/RianNegreiros/go-graphql-api/config"
	"github.com/RianNegreiros/go-graphql-api/internal/db"
	"log"
)

func main() {
	ctx := context.Background()

	config := config.New()

	db := db.New(ctx, config)

	if err := db.Migrate(); err != nil {
		log.Fatalf("error migrating db: %v", err)
	}
}
