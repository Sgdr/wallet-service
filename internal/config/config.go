package config

import (
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Db       DB     `mapstructure:"DB"`
	HttpPort string `mapstructure:"HTTP_PORT"`
}

type DB struct {
	Name           string `mapstructure:"NAME"`
	Host           string `mapstructure:"HOST"`
	Port           string `mapstructure:"PORT"`
	User           string `mapstructure:"USER"`
	Password       string `mapstructure:"PASSWORD"`
	MaxConnections int    `mapstructure:"MAX_CONNECTIONS"`
}

// Load set configuration parameters.
// At first read config from file if configPath parameter is not empty.
// After that read environment variables
func Load(configPath string) (*Config, error) {
	cfg := new(Config)
	if len(configPath) > 0 {
		// read config from file - it will be default values
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
		if err := viper.Unmarshal(cfg); err != nil {
			return nil, err
		}
	}

	// read parameters from environment variables -> they override default values from file
	if envVar, ok := os.LookupEnv("HTTP_PORT"); ok {
		cfg.HttpPort = envVar
	}
	if envVar, ok := os.LookupEnv("DB_USER"); ok {
		cfg.Db.User = envVar
	}
	if envVar, ok := os.LookupEnv("DB_PASSWORD"); ok {
		cfg.Db.Password = envVar
	}
	if envVar, ok := os.LookupEnv("DB_HOST"); ok {
		cfg.Db.Host = envVar
	}
	if envVar, ok := os.LookupEnv("DB_PORT"); ok {
		cfg.Db.Port = envVar
	}
	if envVar, ok := os.LookupEnv("DB_NAME"); ok {
		cfg.Db.Name = envVar
	}
	if envVar, ok := os.LookupEnv("DB_MAX_CONNECTIONS"); ok {
		mc, err := strconv.Atoi(envVar)
		if err != nil {
			return nil, err
		}
		cfg.Db.MaxConnections = mc
	}
	return cfg, nil
}
