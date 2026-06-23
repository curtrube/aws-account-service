package db

import (
	"aws-account-service/internal/account"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type AccountStore struct {
	db DB
}

func NewAccountStore(db DB) *AccountStore {
	return &AccountStore{db: db}
}

// TODO: figure out what the func param is needed to create accounts
//func (s *AccountStore) Create(ctx context.Context) (*account.Account, error) {
//	rows, err := s.db.Query(ctx,
//		`INSERT INTO aws_service.account (
//			id, name, email, joined_method, joined_timestamp, state
//		 )
//		 VALUES ($1, $2, $3, $4, $5, $6)
//		 RETURNING id, name, email, state`,
//		account,
//	)
//	if err != nil {
//		return nil, fmt.Errorf("create account: %w", err)
//	}
//	acct, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[account.Account])
//	if err != nil {
//		return nil, fmt.Errorf("create account: %w", err)
//	}
//	return &acct, nil
//}

func (s *AccountStore) GetByName(ctx context.Context, name string) (*account.Account, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, name, email, joined_method, joined_timestamp, state
	     FROM account
	     WHERE name = $1`,
		name,
	)
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}
	account, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[account.Account])
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}
	return &account, nil
}

func (s *AccountStore) List(ctx context.Context) ([]account.Account, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, name, email, joined_method, joined_timestamp, state
	     FROM aws_service.account
	     ORDER BY name DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[account.Account])
}
