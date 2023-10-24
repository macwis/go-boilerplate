package config

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultLogLevel = logrus.InfoLevel
	serviceName     = "myapp"
)

type Config struct {
	ServiceName    string `mapstructure:"SERVICE_NAME"`
	ServiceVersion string `mapstructure:"SERVICE_VERSION"`
	DatastoreURL   string `mapstructure:"DATASTORE_URL"`
	HTTPPort       string `mapstructure:"HTTP_PORT"`
	GRPCPort       string `mapstructure:"GRPC_PORT"`
	HealthPort     string `mapstructure:"HEALTH_PORT"`

	LogLevel string `mapstructure:"LOG_LEVEL"`
}

func LoadConfig(log *logrus.Logger, configObject *Config, fileNames ...string) (*viper.Viper, error) {
	mainConfig := viper.New()
	fileNames = append([]string{"default.yaml", "config/default.yaml"}, fileNames...)

	for _, fileName := range fileNames {
		viperConfig := viper.New()
		viperConfig.SetConfigFile(fileName)

		if strings.Contains(fileName, "default.") {
			viper.AutomaticEnv()
		}

		if err := viperConfig.MergeInConfig(); err != nil {
			log.WithError(err).Info("config not found; skipping")
			continue
		}

		if err := mainConfig.MergeConfigMap(viperConfig.AllSettings()); err != nil {
			return nil, err
		}
	}

	if err := mainConfig.Unmarshal(configObject); err != nil {
		return nil, errors.Wrap(err, "config parsing error")
	}

	return mainConfig, nil
}

func (c *Config) LogFields() map[string]interface{} {
	return map[string]interface{}{
		"Config.ServiceName":    c.ServiceName,
		"Config.ServiceVersion": c.ServiceVersion,
		"Config.DatastoreURL":   c.DatastoreURL,
		"Config.HTTPPort":       c.HTTPPort,
		"Config.GRPCPort":       c.GRPCPort,
		"Config.HealthPort":     c.HealthPort,
		"Config.LogLevel":       c.LogLevel,
	}
}

func (c *Config) IsValid() error {
	if _, err := logrus.ParseLevel(c.LogLevel); err != nil {
		return err
	}

	// TODO: many other validator function calls

	return nil
}

func (c *Config) GetLogLevel(defaultLevel logrus.Level) logrus.Level {
	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		level = defaultLevel
	}

	return level
}

func newConfig(configFiles []string) (*Config, error) {
	var (
		tmpLog = logrus.New()
		cfg    Config
	)
	tmpLog.SetOutput(os.Stdout)
	tmpLog.SetLevel(defaultLogLevel)

	_, err := LoadConfig(tmpLog, &cfg, configFiles...)
	if err != nil {
		return nil, err
	}

	err = cfg.IsValid()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func NewConfig() (*Config, error) {
	configFiles := []string{
		"default.yaml",
		"config/default.yaml",
		"/config/config.yaml",
		"/vault/secrets/config.yaml",
		".env",
	}

	return newConfig(configFiles)
}
