package db_test

import (
	"aws-account-service/internal/db"
	"aws-account-service/internal/model"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	ctx      = context.Background()
	err      error
	testPool *pgxpool.Pool
)

func setupPGTestContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	img := "postgres:18-alpine"
	dbName := "pg_test"
	dbUser := "postgres"
	dbPassword := "postgres"

	ctr, err := postgres.Run(
		ctx,
		img,
		postgres.WithInitScripts(filepath.Join("..", "..", "sql", "00schema.sql")),
		postgres.WithInitScripts(filepath.Join("..", "..", "testdata", "seed.sql")),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.WithDatabase(dbName),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, fmt.Errorf("run postgres test container: %w", err)
	}
	return ctr, nil
}

func TestMain(m *testing.M) {
	ctr, err := setupPGTestContainer(ctx)
	if err != nil {
		log.Fatalf("setup container: %v", err)
	}
	defer ctr.Terminate(ctx)

	connStr, err := ctr.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("connection string: %v", err)
	}

	testPool, err = db.New(ctx, connStr)
	if err != nil {
		log.Fatalf("new db pool: %v", err)
	}
	defer testPool.Close()

	code := m.Run()
	testPool.Close()
	_ = ctr.Terminate(ctx)

	os.Exit(code)
}

//func TestNew(t *testing.T) {
//	ctx := context.Background()
//
//	pool, err := db.New(ctx, connStr)
//	require.NoError(t, err)
//	defer pool.Close()
//}

func TestAccountStore_Create(t *testing.T) {
	require.NotNil(t, testPool)

	store := db.NewAccountStore(testPool)
	require.NotNil(t, store)

	in := model.Account{
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

func TestAccountStore_GetByName(t *testing.T) {
	require.NotNil(t, testPool)

	store := db.NewAccountStore(testPool)
	require.NotNil(t, store)

	name := "test-account-1"
	acct, err := store.GetByName(ctx, name)
	require.NoError(t, err)

	assert.Equal(t, name, acct.Name)
}

func TestAccountStore_List(t *testing.T) {
	require.NotNil(t, testPool)

	store := db.NewAccountStore(testPool)
	require.NotNil(t, store)

	accounts, err := store.List(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
}
