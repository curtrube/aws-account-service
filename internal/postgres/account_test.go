package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/curtrube/aws-account-service/internal/account"
	"github.com/curtrube/aws-account-service/internal/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountStore_Create(t *testing.T) {
	require.NotNil(t, testPool)

	conn, err := testPool.Acquire(ctx)
	defer conn.Release()
	require.NoError(t, err)

	store := postgres.NewAccountStore(conn)
	require.NotNil(t, store)

	in := account.Account{
		ID:              "333333333333",
		Name:            "test-account-3",
		Email:           "user@user.com",
		JoinedMethod:    "CREATED",
		JoinedTimestamp: time.Now(),
		State:           "ACTIVE",
	}

	out, err := store.Create(ctx, in)
	require.NoError(t, err)

	assert.Equal(t, in.ID, out.ID)
	assert.Equal(t, in.Name, out.Name)
	assert.Equal(t, in.Email, out.Email)
}

func TestAccountStore_UpdateExisting(t *testing.T) {
	require.NotNil(t, testPool)

	conn, err := testPool.Acquire(ctx)
	require.NoError(t, err)
	defer conn.Release()

	in := account.Account{
		ID:              "333333333333",
		Name:            "test-account-3",
		Email:           "user@user.com",
		JoinedMethod:    "CREATED",
		JoinedTimestamp: time.Now(),
		State:           "ACTIVE",
	}

	store := postgres.NewAccountStore(conn)
	require.NotNil(t, store)

	_, err = store.Create(ctx, in)
	require.NoError(t, err)

	in = account.Account{
		ID:              "333333333333",
		Name:            "test-account-3",
		Email:           "user@user.com",
		JoinedMethod:    "CREATED",
		JoinedTimestamp: time.Now(),
		State:           "CLOSED",
	}

	out, err := store.Create(ctx, in)
	require.NoError(t, err)
	assert.Equal(t, "CLOSED", out.State)
}

func TestAccountStore_GetByName(t *testing.T) {
	require.NotNil(t, testPool)

	conn, err := testPool.Acquire(ctx)
	require.NoError(t, err)
	defer conn.Release()

	store := postgres.NewAccountStore(conn)
	require.NotNil(t, store)

	name := "test-account-1"
	acct, err := store.GetByName(ctx, name)
	require.NoError(t, err)

	assert.Equal(t, name, acct.Name)
}

func TestAccountStore_List(t *testing.T) {
	require.NotNil(t, testPool)

	conn, err := testPool.Acquire(ctx)
	defer conn.Release()
	require.NoError(t, err)

	store := postgres.NewAccountStore(conn)
	require.NotNil(t, store)

	accounts, err := store.List(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
}
