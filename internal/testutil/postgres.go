package testutil

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/curtrube/aws-account-service/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupPGTestContainer(ctx context.Context) (*pg.PostgresContainer, error) {
	img := "postgres:18-alpine"
	dbName := "pg_test"
	dbUser := "postgres"
	dbPassword := "postgres"

	ctr, err := pg.Run(
		ctx,
		img,
		pg.WithInitScripts(filepath.Join("..", "..", "sql", "00schema.sql")),
		pg.WithInitScripts(filepath.Join("..", "..", "testdata", "seed.sql")),
		pg.WithUsername(dbUser),
		pg.WithPassword(dbPassword),
		pg.WithDatabase(dbName),
		pg.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, fmt.Errorf("run postgres test container: %w", err)
	}
	return ctr, nil
}

func New(ctx context.Context, t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctr, err := setupPGTestContainer(ctx)
	require.NoErrorf(t, err, "setup test container: %v", err)

	t.Cleanup(func() {
		ctr.Terminate(ctx)
	})

	connStr, err := ctr.ConnectionString(ctx)
	require.NoErrorf(t, err, "get container connection string: %v", err)

	testPool, err := db.New(ctx, connStr)
	require.NoErrorf(t, err, "create new db pool: %v", err)

	t.Cleanup(func() {
		testPool.Close()
	})

	return testPool
}

// Start sets up a postgres test container. Returns a pgxPool, a cleanup func, and an error.
func Start(ctx context.Context) (*pgxpool.Pool, func(), error) {
	ctr, err := setupPGTestContainer(ctx)
	if err != nil {
		return nil, nil, err
	}

	connStr, err := ctr.ConnectionString(ctx)
	if err != nil {
		_ = ctr.Terminate(ctx)
		return nil, nil, err
	}

	pool, err := db.New(ctx, connStr)
	if err != nil {
		_ = ctr.Terminate(ctx)
		return nil, nil, err
	}

	cleanup := func() {
		pool.Close()
		_ = ctr.Terminate(ctx)
	}

	return pool, cleanup, nil
}

//func TestNew(t *testing.T) {
//	ctx := context.Background()
//
//	pool, err := db.New(ctx, connStr)
//	require.NoError(t, err)
//	defer pool.Close()
//}
