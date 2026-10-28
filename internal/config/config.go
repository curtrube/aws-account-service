package config

import (
	"errors"
	"os"
	"strconv"
)

// Config holds cli configuration for aws-account-service.
type Config struct {
	AwsRegion   string
	DatabaseURL string
	Environment string
	Port        int
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
		return nil, errors.New("missing required env variable: DATABASE_URL")
	}

	env, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		return nil, errors.New("missing required env variable: ENVIRONMENT")
	}

	//logLevel, ok := os.LookupEnv("LOG_LEVEL")
	//if !ok {
	//	logLevel = slog.LevelInfo
	//}

	var port int
	var err error
	p, ok := os.LookupEnv("PORT")
	if !ok {
		port = 8080
	} else {
		port, err = strconv.Atoi(p)
		if err != nil {
			return nil, errors.New("parsing port string to integer")
		}
	}

	return &Config{
		AwsRegion:   awsRegion,
		DatabaseURL: databaseURL,
		Environment: env,
		Port:        port,
	}, nil
}
