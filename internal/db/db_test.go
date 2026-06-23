package db_test

import (
	"aws-account-service/internal/db"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupPGTestContainer(ctx context.Context) (string, error) {
	img := "postgres:18-alpine"
	dbName := "pg_test"
	dbUser := "postgres"
	dbPassword := "postgres"

	container, err := postgres.Run(
		ctx,
		img,
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.WithDatabase(dbName),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return "", fmt.Errorf("run postgres test container: %w", err)
	}
	connStr, err := container.ConnectionString(ctx)
	return connStr, nil

	//db, err := pgx.Connect(ctx, connStr)
	//if err != nil {
	//	return nil, nil, err
	//}

	//if err := db.Ping(ctx); err != nil {
	//	return nil, nil, err
	//}

	//return db, container, nil
}

func TestNew(t *testing.T) {
	ctx := context.Background()

	testDatabaseURL, err := setupPGTestContainer(ctx)
	require.NoError(t, err)

	pool, err := db.New(ctx, testDatabaseURL)
	require.NoError(t, err)
	pool.Close()
}
