package main

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	level, err := zap.ParseAtomicLevel(viper.GetString("LoggerLevel"))
	if err != nil {
		return nil, err
	}
	cfg.Level = level
	cfg.OutputPaths = []string{"stdout"}
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	logger, err := cfg.Build(zap.WithCaller(false))
	return logger, err
}
