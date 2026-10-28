// Package account provides the account domain model
// and repository contract.
package account

import (
	"context"
	"time"
)

type Repository interface {
	Upsert(ctx context.Context, account Account) (Account, error)
	List(ctx context.Context) ([]Account, error)
}

// Account represents an AWS Organizations Account.
type Account struct {
	ID           int64     `db:"id"`
	AccountID    string    `db:"account_id"`
	Name         string    `db:"name"`
	Arn          string    `db:"arn"`
	Email        string    `db:"email"`
	State        string    `db:"state"`
	JoinedMethod string    `db:"joined_method"`
	JoinedAt     time.Time `db:"joined_at"`
}
