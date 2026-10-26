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
	account "github.com/curtrube/aws-account-service/internal/aws"
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

	// start the server in the background
	go func() {
		slog.Info("server starting", "address", srv.Addr, "environment", os.Getenv("ENVIRONMENT"))

		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			cancel()
		}
	}()

	var wg sync.WaitGroup

	wg.Go(func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		poll := func() {
			slog.Info("polling AWS account metadata")
			pollCtx, pollCancel := context.WithTimeout(ctx, 30*time.Second)
			defer pollCancel()

			accounts, err := account.ListAccounts(pollCtx, orgClient)
			if err != nil {
				if ctx.Err() == nil {
					slog.Error("listing AWS accounts", "error", err)
				}
				return
			}
			slog.Info("AWS accounts retrieved", "count", len(accounts))
			// Process accounts here.
		}

		for {
			select {
			case <-ctx.Done():
				slog.Info("AWS account poller shutting down")
				return

			case <-ticker.C:
				poll()
			}
		}
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
