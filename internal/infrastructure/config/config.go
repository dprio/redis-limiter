package config

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type (
	Config struct {
		DB          *DB          `mapstructure:"db"`
		Web         *Web         `mapstructure:"web"`
		RateLimiter *RateLimiter `mapstructure:"rate-limiter"`
		Redis       *Redis       `mapstructure:"redis"`
	}

	DB struct {
		Driver   string `mapstructure:"driver"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
	}

	Web struct {
		Port string `mapstructure:"port"`
	}

	RateLimiter struct {
		WindowDuration time.Duration `mapstructure:"window-duration"`
		IPMaxRequests  int64         `mapstructure:"ip-max-requests"`
		BlockDuration  time.Duration `mapstructure:"block-duration"`
	}

	Redis struct {
		Host     string `mapstructure:"host"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	}
)

func New() *Config {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var config Config
	err := viper.Unmarshal(&config, viper.DecodeHook(
		mapstructure.StringToTimeDurationHookFunc(),
	))
	if err != nil {
		panic(err)
	}

	return &config
}
