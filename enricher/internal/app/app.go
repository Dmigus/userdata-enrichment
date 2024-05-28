package app

import (
	"context"
	"enricher/internal/providers/handlingrunner"
	"enricher/internal/providers/messagehandler"
	"enricher/internal/providers/repository"
	"enricher/internal/service"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var Module = fx.Module("app",
	fx.Provide(
		repository.New,
		setupServiceLifecycle,
		messagehandler.New,
		handlingrunner.NewKafkaConsumerGroupRunner,
		newConsumerGroupRunnerConfig,
		repoConnection,
	),
)

type serviceParams struct {
	fx.In
	runner   service.FioHandlingRunner
	enricher *messagehandler.Enricher
	logger   *zap.Logger
	repo     service.Repository
}

func setupServiceLifecycle(lc fx.Lifecycle, params serviceParams) *service.EnrichService {
	s := service.NewEnrichService(params.runner, params.enricher, params.logger, params.repo)
	var cancelCtxFn context.CancelFunc
	sErr := make(chan error)
	lc.Append(fx.StartHook(func(ctx context.Context) {
		var childCtx context.Context
		childCtx, cancelCtxFn = context.WithCancel(ctx)
		go func() {
			defer close(sErr)
			sErr <- s.Run(childCtx)
		}()
	}))
	lc.Append(fx.StopHook(func() error {
		cancelCtxFn()
		return <-sErr
	}))
	return s
}

func repoConnection(config *Config) (*gorm.DB, error) {
	return nil, nil
}

func newConsumerGroupRunnerConfig(config *Config, logger *zap.Logger) handlingrunner.ConsumerGroupRunnerConfig {
	return handlingrunner.ConsumerGroupRunnerConfig{Brokers: config.DataBus.Brokers,
		Topic:  config.DataBus.Topic,
		Logger: logger,
	}
}
