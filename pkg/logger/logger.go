package logger

import (
	"os"

	"go.uber.org/zap"
)

var (
	Logger *zap.SugaredLogger
)

func InitLogger() error {
	var options []zap.Option
	if _, found := os.LookupEnv("HOMECTL_LOGS_VERBOSE"); found {
		options = append(options, zap.IncreaseLevel(zap.DebugLevel))
	}

	logger, err := zap.NewProduction(options...)
	if err != nil {
		return err
	}

	Logger = logger.Sugar()

	return nil
}
