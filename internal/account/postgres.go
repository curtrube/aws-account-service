package account

import (
	"context"
	"fmt"

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
type AccountRepository struct {
	db *pgxpool.Conn
}

// NewAccountRepository creates a new account repository.
func NewAccountRepository(db *pgxpool.Conn) *AccountRepository {
	return &AccountRepository{db: db}
}

const createAccountSQL = `
INSERT INTO aws_cloud.account (
	account_id,
	account_name,
	account_arn,
	email,
	account_state,
	joined_method,
	joined_at
)
VALUES (
	@account_id,
	@account_name,
	@account_arn,
	@email,
	@account_state,
	@joined_method,
	@joined_at
)
ON CONFLICT(id)
DO UPDATE SET
	account_name  = EXCLUDED.account_name,
	email         = EXCLUDED.email,
	account_state = EXCLUDED.account_state
RETURNING 
	id,
	account_id,
	account_name,
	account_arn,
	email,
	account_state,
	joined_method,
	joined_at`

// Create inserts and returns the created account.
// If conflict on (id) will perform an upsert on the row.
func (s *AccountRepository) Create(ctx context.Context, acct Account) (Account, error) {
	rows, err := s.db.Query(
		ctx,
		createAccountSQL,
		pgx.NamedArgs{
			"account_id":    acct.AccountID,
			"account_name":  acct.AccountName,
			"account_arn":   acct.AccountArn,
			"email":         acct.Email,
			"account_state": acct.AccountState,
			"joined_method": acct.JoinedMethod,
			"joined_at":     acct.JoinedAt,
		})
	if err != nil {
		return Account{}, fmt.Errorf("create account: %w", err)
	}
	account, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Account])
	if err != nil {
		return Account{}, fmt.Errorf("create account: %w", err)
	}
	return account, nil
}

// GetByID
// func (s *AccountRepository) GetByID(ctx context.Context) (*account.Account, error) { }

const getAccountSQL = `
SELECT 
	id,
	account_id,
	account_name,
	account_arn,
	email,
	account_state,
	joined_method,
	joined_at
FROM aws_cloud.account
WHERE name = @account_name
`

// GetByName returns an account by its name.
func (s *AccountRepository) GetByName(ctx context.Context, name string) (Account, error) {
	rows, err := s.db.Query(
		ctx,
		getAccountSQL,
		pgx.NamedArgs{
			"account_name": name,
		})
	if err != nil {
		return Account{}, fmt.Errorf("query get account by name: %s: %w", name, err)
	}

	account, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Account])
	if err != nil {
		return Account{}, fmt.Errorf("account name: %s not found: %w", name, err)
	}
	return account, nil
}

//const listAccountsSQL = `
//SELECT
//	id,
//	name,
//	email,
//	joined_method,
//	joined_timestamp,
//	state
//FROM aws.account
//ORDER BY name DESC`
//
//// List all accounts in the store.
//func (s *AccountRepository) List(ctx context.Context) ([]Account, error) {
//	rows, err := s.db.Query(
//		ctx,
//		listAccountsSQL,
//	)
//	if err != nil {
//		return nil, fmt.Errorf("query list accounts: %w", err)
//	}
//	return pgx.CollectRows(rows, pgx.RowToStructByName[Account])
//}
