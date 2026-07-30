package config

import (
	"flag"
	"fmt"
	"log/slog"
)

// Config holds cli configuration for aws-account-service.
type Config struct {
	LogLevel slog.Level
}

// ParseFlags parses command-line arguments into a config.
func ParseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("aws-account-service", flag.ContinueOnError)
	cfg := &Config{}

	var level string
	fs.StringVar(&level, "log-level", "INFO", "Log level (TRACE, DEBUG, INFO, WARN, ERROR)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if err := cfg.LogLevel.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", level, err)
	}

	return cfg, nil
}
