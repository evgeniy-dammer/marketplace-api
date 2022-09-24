package config

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	DBHost string `mapstructure:"DATABASE_HOST"`
	DBUser string `mapstructure:"DATABASE_USER"`
	DBPass string `mapstructure:"DATABASE_PASS"`
	DBName string `mapstructure:"DATABASE_NAME"`
	DBPort string `mapstructure:"DATABASE_PORT"`
}

func LoadConfiguration() (c Configuration, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)

	return
}
