package models

//DbConfiguration model
type DbConfiguration struct {
	DbHost string `mapstructure:"DATABASE_HOST"`
	DbUser string `mapstructure:"DATABASE_USER"`
	DbPass string `mapstructure:"DATABASE_PASS"`
	DbName string `mapstructure:"DATABASE_NAME"`
	DbPort string `mapstructure:"DATABASE_PORT"`
}
