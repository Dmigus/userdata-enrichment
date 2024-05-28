// Package main содержит main для запуска программы, а также обработку переменных окружения и аргументов командной строки
package main

import (
	"enricher/internal/app"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type globResult struct {
	fx.Out
	config *app.Config
	logger *zap.Logger
}

func initGlobalModule(lc fx.Lifecycle) (globResult, error) {
	if err := initViper(); err != nil {
		return globResult{}, err
	}
	config, err := readConfig()
	if err != nil {
		return globResult{}, err
	}
	logger, err := getLogger()
	if err != nil {
		return globResult{}, err
	}
	lc.Append(fx.StopHook(logger.Sync))
	return globResult{logger: logger, config: config}, nil
}

func main() {

	fx.New(fx.Provide(initGlobalModule),
		app.Module,
	).Run()
}
