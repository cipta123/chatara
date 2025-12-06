package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	Environment  string
	MediaStorage string
}

var AppConfig *Config

func Load() error {
	// Load .env file if exists (for local development)
	_ = godotenv.Load()

	AppConfig = &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/umess?sslmode=disable"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		MediaStorage: getEnv("MEDIA_STORAGE", "./uploads"),
	}

	if AppConfig.JWTSecret == "your-secret-key-change-in-production" && AppConfig.Environment == "production" {
		return fmt.Errorf("JWT_SECRET must be set in production")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}


