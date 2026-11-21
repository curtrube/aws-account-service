package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds configuration for aws-account-service.
type Config struct {
	AwsRegion   string
	Environment string
	Port        int
	DB          PostgreSQLConfig
}

// PostgresSQLConfig is Postgresql database configuration.
type PostgreSQLConfig struct {
	Host     string
	Port     int
	Database string
	// Fine for development but will want to move this when we deploy.
	Username string
	Password string
}

// Load reads service configuration from environment varaibles and applies defaults.
func Load() (Config, error) {
	cfg := Config{
		AwsRegion:   getEnv("AWS_REGION", "us-east-1"),
		Environment: getEnv("ENVIRONMENT", "dev"),
		Port:        getInt("PORT", 8080),
		DB: PostgreSQLConfig{
			Host:     mustGetEnv("PGHOST"),
			Port:     getInt("PGPORT", 5432),
			Database: mustGetEnv("PGDATABASE"),
			Username: mustGetEnv("PGUSER"),
			Password: mustGetEnv("PGPASSWORD"),
		},
	}

	return cfg, nil
}

// getEnv retrieves the value of the environment variable identified by key.
// If the variable is unset or empty, it returns fallback.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// mustGetEnv retrieves the value of the environment variable identified by key.
// It panics if the variable is unset or empty.
func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("environment variable %q not set", key))
	}
	return value
}

// getInt retrieves the int value of the environment variable identified by key.
// If the variable is unset or empty, it returns fallback.
func getInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return value
}
