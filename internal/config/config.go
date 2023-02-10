package config

import "github.com/spf13/viper"

// InitConfig initializes configuration
func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	
	return viper.ReadInConfig()
}
