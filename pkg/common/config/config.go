package config

import (
	"github.com/evgeniy-dammer/emenu-api/pkg/common/models"
	"github.com/spf13/viper"
)

// LoadConfiguration loads database configuration from .env file
func LoadConfiguration() (dbConfiguration models.DbConfiguration, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&dbConfiguration)

	return
}
