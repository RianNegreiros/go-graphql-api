package config

import (
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

type database struct {
	URL string
}

type jwt struct {
	Secret string
	Issuer string
}

type env struct {
	BuildEnv string
}

type Config struct {
	Database database
	JWT      jwt
	Env      env
}

func LoadEnv(fileName string) {
	re := regexp.MustCompile(`^(.*` + "go-graphql-api" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + "/" + fileName)
	if err != nil {
		godotenv.Load()
	}
}

func New() *Config {
	return &Config{
		Database: database{
			URL: os.Getenv("DATABASE_URL"),
		},
		JWT: jwt{
			Secret: os.Getenv("JWT_SECRET"),
			Issuer: os.Getenv("DOMAIN"),
		},
		Env: env{
			BuildEnv: os.Getenv("BUILD_ENV"),
		},
	}
}
