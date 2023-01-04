package config

import (
	"github.com/evgeniy-dammer/emenu-api/pkg/common/models"
	"github.com/spf13/viper"
)

// LoadConfiguration loads configuration from .env file
func LoadConfiguration() (configuration models.Configuration, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&configuration)

	return
}
