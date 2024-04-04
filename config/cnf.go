package config

import (
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const DefaultFile = ".env"

type Config struct {
	AppConfig         AppConfig
	DbConfig          DbConfig
	MigrationPath     string
	ExternalAPIConfig ExternalAPIConfig
}

func LoadDefault() (*Config, error) {
	return Load(DefaultFile)
}

func Load(file string) (*Config, error) {
	if err := setupViper(file); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	c := &Config{
		AppConfig:         getAppConfig(),
		DbConfig:          getDbConfig(),
		MigrationPath:     getString("MIGRATION_PATH"),
		ExternalAPIConfig: getExternalAPIConfig(),
	}

	return c, nil
}

func setupViper(file string) error {
	viper.SetConfigFile(file)
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func getString(key string) string {
	if !viper.IsSet(key) {
		log.Fatal().
			Str("key", key).
			Msg("Unable to find config value for key")
	}

	return viper.GetString(key)
}

func getInt(key string) int {
	if !viper.IsSet(key) {
		log.Fatal().
			Str("key", key).
			Msg("Unable to find config value for key")
	}
	val := viper.GetInt(key)

	return val
}

func getIntOrDefault(key string, def int) int {
	if !viper.IsSet(key) {
		return def
	}

	return viper.GetInt(key)
}
