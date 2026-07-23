package configs

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3/log"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort          string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPass              string
	DBName              string
	DBSSLMode           string
	Environment         string
	Version             string
	JWTSecret           string
	JWTRefreshSecret    string
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Info("No .env file found, relying on environment variables")
	}
	cfg := &Config{
		ServerPort:          GetEnv("SERVER_PORT", "8080"),
		DBHost:              GetEnv("DB_HOST", "localhost"),
		DBPort:              GetEnv("DB_PORT", "5432"),
		DBUser:              GetEnv("DB_USER", "postgres"),
		DBPass:              GetEnv("DB_PASS", "password"),
		DBName:              GetEnv("DB_NAME", "dbname"),
		DBSSLMode:           GetEnv("DB_SSLMODE", "disable"),
		Environment:         GetEnv("ENVIRONMENT", "development"),
		Version:             GetEnv("VERSION", "1.0.0"),
		JWTSecret:           GetEnv("JWT_SECRET", "your-default-jwt-secret"),
		JWTRefreshSecret:    GetEnv("JWT_REFRESH_SECRET", "your-default-jwt-refresh-secret"),
		CloudinaryCloudName: GetEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:    GetEnv("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret: GetEnv("CLOUDINARY_API_SECRET", ""),
	}

	// Keamanan: di production, secret JWT default = siapa pun bisa memalsukan token admin.
	if cfg.Environment == "production" {
		if cfg.JWTSecret == "your-default-jwt-secret" || cfg.JWTRefreshSecret == "your-default-jwt-refresh-secret" {
			log.Fatal("JWT_SECRET dan JWT_REFRESH_SECRET wajib diisi (bukan nilai default) di production")
		}
	}

	return cfg, nil
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName, c.DBSSLMode, "UTC",
	)
}
