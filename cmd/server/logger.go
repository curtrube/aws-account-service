package main

import (
	"log/slog"
	"os"
	"strings"
)

func logger() *slog.Logger {
	level := parseLevel(os.Getenv("LOG_LEVEL"))

	opts := slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if strings.EqualFold(os.Getenv("LOG_FORMAT"), "text") {
		handler = slog.NewTextHandler(os.Stdout, &opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &opts)
	}

	return slog.New(handler)
}

// parseLevel parses a string log level into a slog.Level.
func parseLevel(value string) slog.Level {
	switch strings.ToLower(value) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
