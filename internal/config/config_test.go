package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, cfg.AwsRegion, "us-east-1")
	assert.Equal(t, cfg.Port, 8080)
}

func TestLoad_Overrides(t *testing.T) {
	databaseURL := "postgresql://devuser:devpassword@localhost:5432/devdb"
	os.Setenv("DATABASE_URL", databaseURL)
	os.Setenv("ENVIRONMENT", "test")

	cfg, err := Load()

	require.NoError(t, err)

	assert.Equal(t, cfg.AwsRegion, "us-east-1")
	assert.Equal(t, cfg.DatabaseURL, databaseURL)
	assert.Equal(t, cfg.Environment, "test")
	assert.Equal(t, cfg.Port, 8080)
}
