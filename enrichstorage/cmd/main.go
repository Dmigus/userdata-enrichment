// Package main содержит main для запуска программы, а также обработку переменных окружения и аргументов командной строки
package main

import (
	"enrichstorage/internal/apps/enrichstorage"
	"enrichstorage/internal/apps/outboxsender"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type globResult struct {
	fx.Out
	EnrichConfig       *enrichstorage.Config
	OutboxSenderConfig *outboxsender.Config
	Logger             *zap.Logger
}

func initGlobalModule(lc fx.Lifecycle) (globResult, error) {
	enrichConfig := enrichstorage.Config{}
	outboxSenderConfig := outboxsender.Config{}
	err := parseConfigs(&enrichConfig, &outboxSenderConfig)
	if err != nil {
		return globResult{}, err
	}
	logger, err := getLogger()
	if err != nil {
		return globResult{}, err
	}
	lc.Append(fx.StopHook(logger.Sync))
	return globResult{Logger: logger, EnrichConfig: &enrichConfig, OutboxSenderConfig: &outboxSenderConfig}, nil
}

func main() {
	fx.New(fx.Provide(initGlobalModule),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		enrichstorage.Module,
		outboxsender.Module,
	).Run()
}
