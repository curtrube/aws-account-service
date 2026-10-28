package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/curtrube/aws-account-service/internal/account"
	"github.com/curtrube/aws-account-service/internal/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	pool, cleanup, err := testutil.Start(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres test database: %v\n", err)
		os.Exit(1)
	}
	testPool = pool
	code := m.Run()
	cleanup()
	os.Exit(code)
}

func TestRepository_Upsert(t *testing.T) {
	ctx := context.Background()
	require.NotNil(t, testPool)

	repo := NewRepository(testPool)
	require.NotNil(t, repo)

	mgmtAccountID := "111111111111"
	ouID := "o-abc123"
	accountID := "333333333333"
	accountArn := fmt.Sprintf("arn:aws:organizations::%s:account/%s/%s", mgmtAccountID, ouID, accountID)

	in := account.Account{
		AccountID:    accountID,
		Name:         "test-account-3",
		Arn:          accountArn,
		Email:        "foobar@baz.com",
		State:        "ACTIVE",
		JoinedMethod: "CREATED",
		JoinedAt:     time.Now(),
	}

	out, err := repo.Upsert(ctx, in)
	require.NoError(t, err)

	assert.NotEmpty(t, out.ID)
	assert.Equal(t, in.AccountID, out.AccountID)
	assert.Equal(t, in.Name, out.Name)
	assert.Equal(t, in.Arn, out.Arn)
	assert.Equal(t, in.Email, out.Email)
	assert.Equal(t, in.State, out.State)
	assert.Equal(t, in.JoinedMethod, out.JoinedMethod)
}

//func TestAccountStore_UpdateExisting(t *testing.T) {
//	require.NotNil(t, testPool)
//
//	conn, err := testPool.Acquire(ctx)
//	require.NoError(t, err)
//	defer conn.Release()
//
//	in := account.Account{
//		ID:              "333333333333",
//		Name:            "test-account-3",
//		Email:           "user@user.com",
//		JoinedMethod:    "CREATED",
//		JoinedTimestamp: time.Now(),
//		State:           "ACTIVE",
//	}
//
//	store := postgres.NewAccountStore(conn)
//	require.NotNil(t, store)
//
//	_, err = store.Create(ctx, in)
//	require.NoError(t, err)
//
//	in = account.Account{
//		ID:              "333333333333",
//		Name:            "test-account-3",
//		Email:           "user@user.com",
//		JoinedMethod:    "CREATED",
//		JoinedTimestamp: time.Now(),
//		State:           "CLOSED",
//	}
//
//	out, err := store.Create(ctx, in)
//	require.NoError(t, err)
//	assert.Equal(t, "CLOSED", out.State)
//}
//
//func TestAccountStore_GetByName(t *testing.T) {
//	require.NotNil(t, testPool)
//
//	conn, err := testPool.Acquire(ctx)
//	require.NoError(t, err)
//	defer conn.Release()
//
//	store := postgres.NewAccountStore(conn)
//	require.NotNil(t, store)
//
//	name := "test-account-1"
//	acct, err := store.GetByName(ctx, name)
//	require.NoError(t, err)
//
//	assert.Equal(t, name, acct.Name)
//}
//
//func TestAccountStore_List(t *testing.T) {
//	require.NotNil(t, testPool)
//
//	conn, err := testPool.Acquire(ctx)
//	defer conn.Release()
//	require.NoError(t, err)
//
//	store := postgres.NewAccountStore(conn)
//	require.NotNil(t, store)
//
//	accounts, err := store.List(context.Background())
//	require.NoError(t, err)
//	require.NotEmpty(t, accounts)
//}
