// Package main содержит main для запуска программы, а также обработку переменных окружения и аргументов командной строки
package main

import (
	"enricher/internal/app"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type globResult struct {
	fx.Out
	Config *app.Config
	Logger *zap.Logger
}

func initGlobalModule(lc fx.Lifecycle) (globResult, error) {
	if err := initViper(); err != nil {
		return globResult{}, err
	}
	config, err := parseConfig()
	if err != nil {
		return globResult{}, err
	}
	logger, err := getLogger()
	if err != nil {
		return globResult{}, err
	}
	lc.Append(fx.StopHook(logger.Sync))
	return globResult{Logger: logger, Config: config}, nil
}

func main() {
	fx.New(fx.Provide(initGlobalModule),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		app.Module,
	).Run()
}
