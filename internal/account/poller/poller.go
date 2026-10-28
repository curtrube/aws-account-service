// Package poller
package poller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/curtrube/aws-account-service/internal/account"
)

type ListAccountsFunc func(ctx context.Context) ([]account.Account, error)

type Poller struct {
	listAccounts ListAccountsFunc
	repository   account.Repository
	interval     time.Duration
}

func New(
	listAccounts ListAccountsFunc,
	repository account.Repository,
	interval time.Duration,
) *Poller {
	return &Poller{
		listAccounts: listAccounts,
		repository:   repository,
		interval:     interval,
	}
}

// Poll lists all accounts and performs an upsert.
func (p *Poller) Poll(ctx context.Context) error {
	slog.Info("poller listing accounts")
	accounts, err := p.listAccounts(ctx)
	if err != nil {
		return fmt.Errorf("list accounts: %w", err)
	}

	for _, a := range accounts {
		if _, err := p.repository.Upsert(ctx, a); err != nil {
			return fmt.Errorf("upsert account %s: %w", a.AccountID, err)
		}
	}
	return nil
}

// Run polls for account updates at the configured interval until ctx is cancelled.
// It waits for the first tick before polling. Individual poll errors are logged,
// and the polling loop continues.
func (p *Poller) Run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("account poller shutting down")
			return

		case <-ticker.C:
			pollCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			err := p.Poll(pollCtx)
			cancel()

			if err != nil && ctx.Err() == nil {
				slog.Error("polling accounts", "error", err)
			}
		}
	}
}
