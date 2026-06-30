package postgres_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/curtrube/aws-account-service/internal/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	ctx      = context.Background()
	testPool *pgxpool.Pool
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

	testPool, err = postgres.New(ctx, connStr)
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
