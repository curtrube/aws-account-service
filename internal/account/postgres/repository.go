// Package postgres provides PostgreSQL persistence adapter.
package postgres

import (
	"context"
	"fmt"

	"github.com/curtrube/aws-account-service/internal/account"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository persists and retrieves account data from PostgreSQL.
type Repository struct {
	db *pgxpool.Pool
}

// compile time check that *Repository implements the account.Repository inteface.
var _ account.Repository = (*Repository)(nil)

// NewRepository creates a new account repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

const upsertAccountSQL = `
INSERT INTO aws_cloud.account (
	account_id,
	name,
	arn,
	email,
	state,
	joined_method,
	joined_at
)
VALUES (
	@account_id,
	@name,
	@arn,
	@email,
	@state,
	@joined_method,
	@joined_at
)
ON CONFLICT(account_id)
DO UPDATE SET
	name          = EXCLUDED.name,
	arn           = EXCLUDED.arn,
	email         = EXCLUDED.email,
	state         = EXCLUDED.state,
	joined_method = EXCLUDED.joined_method,
	joined_at     = EXCLUDED.joined_at
RETURNING 
	id,
	account_id,
	name,
	arn,
	email,
	state,
	joined_method,
	joined_at`

// Upsert inserts a new row into the account table or updates an existing row.
func (r *Repository) Upsert(ctx context.Context, acct account.Account) (account.Account, error) {
	var saved account.Account
	err := r.db.QueryRow(
		ctx,
		upsertAccountSQL,
		pgx.NamedArgs{
			"account_id":    acct.AccountID,
			"name":          acct.Name,
			"arn":           acct.Arn,
			"email":         acct.Email,
			"state":         acct.State,
			"joined_method": acct.JoinedMethod,
			"joined_at":     acct.JoinedAt,
		}).Scan(
		&saved.ID,
		&saved.AccountID,
		&saved.Name,
		&saved.Arn,
		&saved.Email,
		&saved.State,
		&saved.JoinedMethod,
		&saved.JoinedAt,
	)
	if err != nil {
		return account.Account{}, fmt.Errorf("upsert account: %w", err)
	}
	return saved, nil
}

// List retrieves all accounts from the account table.
func (r *Repository) List(ctx context.Context) ([]account.Account, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT id, account_id, name, arn, email, state, joined_method, joined_at
		 FROM aws_cloud.account`,
	)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	defer rows.Close()

	var out []account.Account

	for rows.Next() {
		a := account.Account{}
		err = rows.Scan(
			&a.ID,
			&a.AccountID,
			&a.Name,
			&a.Arn,
			&a.Email,
			&a.State,
			&a.JoinedMethod,
			&a.JoinedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("list accounts scan rows: %w", err)
		}
		out = append(out, a)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return out, nil
}

// GetByID
// func (s *AccountRepository) GetByID(ctx context.Context) (*account.Account, error) { }

//const getAccountByNameSQL = `
//SELECT
//	id,
//	account_id,
//	name,
//	arn,
//	email,
//	state,
//	joined_method,
//	joined_at
//FROM aws_cloud.account
//WHERE name = @name`
//
//// GetByName returns an account by that name.
//func (r *Repository) GetByName(ctx context.Context, name string) (account.Account, error) {
//	rows, err := r.db.Query(
//		ctx,
//		getAccountByNameSQL,
//		pgx.NamedArgs{
//			"name": name,
//		})
//	if err != nil {
//		return account.Account{}, fmt.Errorf("query get account by name: %s: %w", name, err)
//	}
//
//	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[account.Account])
//	if err != nil {
//		return account.Account{}, fmt.Errorf("account name: %s not found: %w", name, err)
//	}
//	return row, nil
//}

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
