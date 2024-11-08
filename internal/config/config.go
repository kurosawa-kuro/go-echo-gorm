package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port        string
	DatabaseURL string
}

func Load() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	if err := godotenv.Load(fmt.Sprintf("configs/.env.%s", env)); err != nil {
		return nil, fmt.Errorf("error loading .env.%s file: %w", env, err)
	}

	return &Config{
		Environment: env,
		Port:        os.Getenv("PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}, nil
}
