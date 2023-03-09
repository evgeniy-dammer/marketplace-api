package logger

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() error {
	var err error

	if viper.GetBool("service.production") {
		Logger, err = zap.NewProduction()

		return err
	}

	Logger, err = zap.NewDevelopment()

	return err
}
