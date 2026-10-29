// Package account provides the account domain model
// and repository contract.
package account

import (
	"context"
	"time"
)

type Repository interface {
	Get(ctx context.Context, id int) (Account, error)
	List(ctx context.Context) ([]Account, error)
	Upsert(ctx context.Context, account Account) (Account, error)
}

// Account represents an AWS Organizations Account.
type Account struct {
	ID           int       `db:"id"`
	AccountID    string    `db:"account_id"`
	Name         string    `db:"name"`
	Arn          string    `db:"arn"`
	Email        string    `db:"email"`
	State        string    `db:"state"`
	JoinedMethod string    `db:"joined_method"`
	JoinedAt     time.Time `db:"joined_at"`
}
