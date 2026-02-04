package configs

import (
	"os"

	"github.com/gofiber/fiber/v3/log"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
	Environment string
	Version     string
	JWTSecret   string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file found, relying on environment variables")
	}
	return &Config{
		ServerPort:  GetEnv("SERVER_PORT", "8080"),
		DatabaseURL: GetEnv("DATABASE_URL", "postgres://user:password@localhost:5432/dbname"),
		Environment: GetEnv("ENVIRONMENT", "development"),
		Version:     GetEnv("VERSION", "1.0.0"),
		JWTSecret:   GetEnv("JWT_SECRET", "your-default-jwt-secret"),
	}, nil
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
