package config

import (
	"github.com/joho/godotenv"
	"os"
)

type database struct {
	URL string
}

type Config struct {
	Database database
}

func New() *Config {
	godotenv.Load()

	return &Config{
		Database: database{
			URL: os.Getenv("DATABASE_URL"),
		},
	}
}
