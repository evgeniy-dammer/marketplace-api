package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// InitConfig initializes configuration.
func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return errors.Wrap(viper.ReadInConfig(), "can not read config")
}
