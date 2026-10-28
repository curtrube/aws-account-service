package config

import (
	"errors"
	"os"
)

// Config holds cli configuration for aws-account-service.
type Config struct {
	AwsRegion   string
	DatabaseURL string
	Environment string
	Port        string
}

// Load reads service configuration from environment varaibles and applies defaults.
func Load() (*Config, error) {
	awsRegion, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		awsRegion = ""
	}

	// If not databaseURL, then require PG_*
	databaseURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return nil, errors.New("m")
	}

	env, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		return nil, errors.New("missing required env variable: ENVIRONMENT")
	}

	//logLevel, ok := os.LookupEnv("LOG_LEVEL")
	//if !ok {
	//	logLevel = slog.LevelInfo
	//}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	return &Config{
		AwsRegion:   awsRegion,
		DatabaseURL: databaseURL,
		Environment: env,
		//LogLevel:    logLevel,
		Port: port,
	}, nil
}
