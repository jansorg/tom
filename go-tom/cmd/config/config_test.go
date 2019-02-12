package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/config"
	"github.com/jansorg/tom/go-tom/test_setup"
)

func TestConfig(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	_, err = doConfigCommand("unknown")
	require.Error(t, err)

	data, err := doConfigCommand("json", config.KeyDataDir)
	require.NoError(t, err)
	require.Contains(t, string(data), viper.GetString(config.KeyDataDir))

	data, err = doConfigCommand("yaml", config.KeyDataDir)
	require.NoError(t, err)
	require.Contains(t, string(data), viper.GetString(config.KeyDataDir))
}
