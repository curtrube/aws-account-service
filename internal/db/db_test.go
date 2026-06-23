package db_test

import (
	"aws-account-service/internal/db"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
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

var testPool *pgxpool.Pool
var err error

func TestMain(m *testing.M) {
	ctx := context.Background()

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

func TestAccountStore_List(t *testing.T) {
	require.NotNil(t, testPool)

	accountStore := db.NewAccountStore(testPool)
	require.NotNil(t, accountStore)

	accounts, err := accountStore.List(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
}
