package config

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlags_Defaults(t *testing.T) {
	cfg, err := ParseFlags([]string{})
	require.NoError(t, err)
	assert.Equal(t, slog.LevelInfo, cfg.LogLevel)
}

func TestParseFlags_Overrides(t *testing.T) {
	cfg, err := ParseFlags([]string{
		"--log-level", "debug",
	})
	require.NoError(t, err)
	assert.Equal(t, slog.LevelDebug, cfg.LogLevel)
}
