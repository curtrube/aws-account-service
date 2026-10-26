package account

import (
	"context"
	"time"
)

type Account struct {
	ID           string    `db:"id"`
	AccountID    string    `db:"account_id"`
	AccountName  string    `db:"account_name"`
	AccountArn   string    `db:"account_arn"`
	Email        string    `db:"email"`
	AccountState string    `db:"account_state"`
	JoinedMethod string    `db:"joined_method"`
	JoinedAt     time.Time `db:"joined_at"`
}

type Repository interface {
	Create(ctx context.Context) (Account, error)
	GetByID(ctx context.Context, id string) (Account, error)
	GetByName(ctx context.Context, name string) (Account, error)
	Update(ctx context.Context) error
	List(ctx context.Context) ([]Account, error)
	Delete(ctx context.Context) error
}
