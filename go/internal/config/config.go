package config

import (
	"fmt"
	"os"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	HTTPPort string
	DBDSN    string
}

// LoadFromEnv initializes configuration values from system environment variables with fallback defaults.
func LoadFromEnv() *Config {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	dbHost := getEnv("DB_HOST", "postgres")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USERNAME", "lgpk_user")
	dbPass := getEnv("DB_PASSWORD", "lgpk_secret")
	dbName := getEnv("DB_DATABASE", "lgpk_db")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	return &Config{
		HTTPPort: port,
		DBDSN:    dsn,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}