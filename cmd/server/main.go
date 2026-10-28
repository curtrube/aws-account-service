package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/curtrube/aws-account-service/internal/account"
	awsorg "github.com/curtrube/aws-account-service/internal/account/aws/organizations"
	"github.com/curtrube/aws-account-service/internal/account/poller"
	accountrepo "github.com/curtrube/aws-account-service/internal/account/postgres"
	"github.com/curtrube/aws-account-service/internal/db"
)

func main() {
	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// Shared cancellation context for background work.
	ctx, cancel := context.WithCancel(signalCtx)

	slog.SetDefault(logger())

	// Initalize a new PG pool.
	databaseURL := "postgresql://devuser:devpassword@localhost:5432/devdb"
	pool, err := db.New(ctx, databaseURL)
	if err != nil {
		slog.Error("creating new pgx pool", "error", err.Error())
		return
	}
	defer pool.Close()

	// Initialize AWS before starting the server.
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		slog.Error("unable to load AWS config", "error", err)
		return
	}
	orgClient := organizations.NewFromConfig(cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Report unexpected server errors to main.
	serverErr := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Go(func() {
		slog.Info("server starting", "address", srv.Addr, "environment", os.Getenv("ENVIRONMENT"))

		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			cancel()
		}
	})

	repo := accountrepo.NewRepository(pool)
	accountPoller := poller.New(
		func(ctx context.Context) ([]account.Account, error) {
			return awsorg.ListAccounts(ctx, orgClient)
		},
		repo,
		10*time.Second,
	)

	wg.Go(func() {
		accountPoller.Run(ctx)
	})

	// Wait for a termination signal or unexpected server failure.
	select {
	case <-signalCtx.Done():
		slog.Info("termination signal received")
	case err := <-serverErr:
		slog.Error("HTTP server failed", "error", err)
	}

	// Stop background work.
	cancel()

	// Stop accepting new HTTP connections and wait for active requests.
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown failed", "error", err)

		if closeErr := srv.Close(); closeErr != nil {
			slog.Error("forcing http server close failed", "error", closeErr)
		}
	}

	// Wait for poller to finish.
	wg.Wait()

	slog.Info("server stopped gracefully")
}
