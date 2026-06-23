package db

import (
	"aws-account-service/internal/model"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DB defines database operations required by the account store.
type DB interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

// AccountStore provides account persistence operations.
type AccountStore struct {
	db DB
}

// NewAccountStore creates a new AccountStore.
func NewAccountStore(db DB) *AccountStore {
	return &AccountStore{db: db}
}

// Create inserts a new account and returns the created account.
func (s *AccountStore) Create(ctx context.Context, account model.Account) (*model.Account, error) {
	rows, err := s.db.Query(ctx,
		`INSERT INTO aws_service.account (
			id, name, email, joined_method, joined_timestamp, state
		)
		VALUES (@id, @name, @email, @joined_method, @joined_timestamp, @state)
		RETURNING id, name, email, joined_method, joined_timestamp, state`,
		pgx.NamedArgs{
			"id": account.ID,
			"name": account.Name,
			"email": account.Email,
			"joined_method": account.JoinedMethod,
			"joined_timestamp": account.JoinedTimestamp,
			"state": account.State,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}
	acct, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Account])
	if err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}
	return &acct, nil
}

// GetByName returns an account by its name.
func (s *AccountStore) GetByName(ctx context.Context, name string) (*model.Account, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, name, email, joined_method, joined_timestamp, state
	     FROM aws_service.account
	     WHERE name = $1`,
		name,
	)
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}

	acct, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Account])
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}
	return &acct, nil
}

// List all accounts in the store.
func (s *AccountStore) List(ctx context.Context) ([]model.Account, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, name, email, joined_method, joined_timestamp, state
	     FROM aws_service.account
	     ORDER BY name DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Account])
}
