package account

import (
	"context"
	"time"
)

type Account struct {
	ID              string    `db:"id"`
	Name            string    `db:"name"`
	Email           string    `db:"email"`
	JoinedMethod    string    `db:"joined_method"`
	JoinedTimestamp time.Time `db:"joined_timestamp"`
	State           string    `db:"state"`
}

type AccountService interface {
	Create(ctx context.Context) (Account, error)
	GetByID(ctx context.Context, id string) (Account, error)
	GetByName(ctx context.Context, name string) (Account, error)
	Update(ctx context.Context) error
	List(ctx context.Context) ([]Account, error)
	Delete(ctx context.Context) error
}
