package postgres

import (
	"context"
	"fmt"

	"github.com/curtrube/aws-account-service/internal/account"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB defines database operations required by the account store.
//type DB interface {
//	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
//	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
//	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
//}

// AccountRepository persists and retrieves account data from PostgreSQL.
type AccountStore struct {
	db *pgxpool.Conn
}

// NewAccountRepository creates a new account repository.
func NewAccountStore(db *pgxpool.Conn) *AccountStore {
	return &AccountStore{db: db}
}

const createAccountSQL = `
INSERT INTO aws.account (
	id,
	name,
	email,
	joined_method,
	joined_timestamp,
	state
)
VALUES (
	@id,
	@name,
	@email,
	@joined_method,
	@joined_timestamp,
	@state
)
ON CONFLICT(id)
DO UPDATE SET
	name = EXCLUDED.name,
	email = EXCLUDED.email,
	joined_method = EXCLUDED.joined_method,
	joined_timestamp = EXCLUDED.joined_timestamp,
	state = EXCLUDED.state
RETURNING 
	id,
	name,
	email,
	joined_method,
	joined_timestamp,
	state
`

// Create inserts and returns the created account.
// If conflict on (id) will perform an upsert on the row.
func (s *AccountStore) Create(ctx context.Context, acct account.Account) (*account.Account, error) {
	rows, err := s.db.Query(ctx, createAccountSQL, pgx.NamedArgs{
		"id":               acct.ID,
		"name":             acct.Name,
		"email":            acct.Email,
		"joined_method":    acct.JoinedMethod,
		"joined_timestamp": acct.JoinedTimestamp,
		"state":            acct.State,
	})
	if err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}
	createdAcct, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[account.Account])
	if err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}
	return &createdAcct, nil
}

// GetByID
// func (s *AccountRepository) GetByID(ctx context.Context) (*account.Account, error) { }

const getAccountSQL = `
SELECT 
	id,
	name,
	email,
	joined_method,
	joined_timestamp,
	state
FROM aws.account
WHERE name = @name
`

// GetByName returns an account by its name.
func (s *AccountStore) GetByName(ctx context.Context, name string) (*account.Account, error) {
	rows, err := s.db.Query(ctx, getAccountSQL, pgx.NamedArgs{
		"name": name,
	})
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}

	acct, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[account.Account])
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}
	return &acct, nil
}

const listAccountsSQL = `
SELECT
	id,
	name,
	email,
	joined_method,
	joined_timestamp,
	state
FROM aws.account
ORDER BY name DESC`

// List all accounts in the store.
func (s *AccountStore) List(ctx context.Context) ([]account.Account, error) {
	rows, err := s.db.Query(ctx, listAccountsSQL)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[account.Account])
}
