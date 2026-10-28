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
	t.Parallel()

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

func TestRepository_List(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	require.NotNil(t, testPool)

	repo := NewRepository(testPool)
	require.NotNil(t, repo)
	//{2 222222222222 test-account-2 arn:aws:organizations::111111111111:account/o-abc123/222222222222 bar@baz.com ACTIVE INVITED 2026-06-23 09:30:00 -0500 CDT}
	//{3 333333333333 test-account-3 arn:aws:organizations::111111111111:account/o-abc123/333333333333 foobar@baz.com ACTIVE CREATED 2026-10-28 14:59:55.733278 -0500 CDT}
	//loc := time.FixedZone("CDT", -5*60*60)
	//expected := []account.Account{
	//	{
	//		ID:           1,
	//		AccountID:    "111111111111",
	//		Name:         "test-account-1",
	//		Arn:          "arn:aws:organizations::111111111111:account/o-abc123/111111111111",
	//		Email:        "foo@bar.com",
	//		State:        "ACTIVE",
	//		JoinedMethod: "CREATED",
	//		JoinedAt:     time.Date(2026, time.June, 23, 9, 30, 0, 0, loc),
	//	},
	//}

	accounts, err := repo.List(ctx)
	require.NoError(t, err)

	assert.Equal(t, 3, len(accounts))
	for _, a := range accounts {
		assert.NotEmpty(t, a.ID)
		assert.NotEmpty(t, a.AccountID)
		assert.NotEmpty(t, a.Name)
		assert.NotEmpty(t, a.Arn)
		assert.NotEmpty(t, a.Email)
		assert.NotEmpty(t, a.State)
		assert.NotEmpty(t, a.JoinedMethod)
		assert.NotEmpty(t, a.JoinedAt)
	}

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
