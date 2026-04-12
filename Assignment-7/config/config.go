package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret string
	Port      string
}

func Load() (*Config, error) {
	_ = godotenv.Load() // не критично, если нет файла

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_secret_key"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	return &Config{
		JWTSecret: jwtSecret,
		Port:      port,
	}, nil
}
