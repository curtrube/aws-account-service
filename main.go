package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/curtrube/aws-account-service/internal/config"
)

func main() {
	cfg, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("parsing config: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))

	logger.Info("server starting...", slog.Any("level", cfg.LogLevel))
}
