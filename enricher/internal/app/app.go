package app

import (
	"context"
	"enricher/internal/providers/handlingrunner"
	"enricher/internal/providers/messagehandler"
	"enricher/internal/providers/messagehandler/computers"
	"enricher/internal/providers/repository"
	"enricher/internal/service"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

var Module = fx.Module("app",
	fx.Provide(
		repository.New,
		setupServiceLifecycle,
		messagehandler.New,
		handlingrunner.NewKafkaConsumerGroupRunner,
		newConsumerGroupRunnerConfig,
		repoConnection,
		computers.NewHttpQueryPerformer,

		fx.Annotate(
			agifyAddress,
			fx.ResultTags(`name:"agifyAddress"'`),
		),
		fx.Annotate(
			computers.NewAgifyComputer,
			fx.ParamTags(`name:"agifyAddress"'`, ``),
		),

		fx.Annotate(
			genderizeAddress,
			fx.ResultTags(`name:"genderizeAddress"'`),
		),
		fx.Annotate(
			computers.NewSexComputer,
			fx.ParamTags(`name:"sexAdgenderizeAddressdress"'`, ``),
		),

		fx.Annotate(
			nationalityAddress,
			fx.ResultTags(`name:"nationalityAddress"'`),
		),
		fx.Annotate(
			computers.NewNationalityComputer,
			fx.ParamTags(`name:"nationalityAddress"'`, ``),
		),
	),
	fx.Supply(http.Client{}),
	fx.Decorate(decorateLogger),
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

func decorateLogger(logger *zap.Logger) *zap.Logger {
	return logger.Named("app")
}

func repoConnection(lc fx.Lifecycle, config *Config) (*gorm.DB, error) {
	dsn := config.Repository.GetPostgresDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqldb, err := db.DB()
	if err != nil {
		return nil, err
	}
	lc.Append(fx.StopHook(func() error {
		return sqldb.Close()
	}))
	return db, nil
}

func newConsumerGroupRunnerConfig(config *Config, logger *zap.Logger) handlingrunner.ConsumerGroupRunnerConfig {
	return handlingrunner.ConsumerGroupRunnerConfig{Brokers: config.DataBus.Brokers,
		Topic:  config.DataBus.Topic,
		Logger: logger,
	}
}

func agifyAddress(config *Config) string {
	return config.AgifyAddress
}

func genderizeAddress(config *Config) string {
	return config.GenderizeAddress
}

func nationalityAddress(config *Config) string {
	return config.NationalityAddress
}
