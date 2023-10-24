package config

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func loadConfig(t *testing.T) *Config {
	var cfg Config
	l := logrus.New()

	if _, err := LoadConfig(
		l,
		&cfg,
		"default.yaml", "config/default.yaml", "/config/config.yaml", ".env"); err != nil {
		t.Errorf("error loading config")
	}

	return &cfg
}

func TestConfig_LoadConfig(t *testing.T) {
	t.Parallel()

	cfg := loadConfig(t)
	require.NotNil(t, cfg)
}

func Test_Config(t *testing.T) {
	t.Parallel()

	cfg := Config{
		ServiceName:    "myapp",
		ServiceVersion: "0.1",
		DatastoreURL:   "",
	}

	t.Run("should log fields without errors", func(t *testing.T) {
		cfg.LogFields()
	})

	t.Run("should be valid configuration", func(t *testing.T) {
		err := cfg.IsValid()
		require.NoError(t, err)
	})

	// TODO: run checks for all mandatory getters within the Config
}
