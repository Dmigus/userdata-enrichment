package app

import (
	"context"
	"enricher/internal/providers/handlingrunner"
	"enricher/internal/providers/messagehandler"
	"enricher/internal/providers/messagehandler/computers"
	"enricher/internal/service"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

var Module = fx.Module("app",
	fx.Provide(
		fx.Annotate(
			storage.New,
			fx.As(new(service.Storage))),
		setupServiceLifecycle,
		messagehandler.New,

		fx.Annotate(
			setupRunner,
			fx.As(new(service.FioHandlingRunner)),
		),
		repoConnection,
		fx.Annotate(
			computers.NewHttpQueryPerformer,
			fx.As(new(computers.CallPerformer)),
		),

		fx.Annotate(
			agifyAddress,
			fx.ResultTags(`name:"agifyAddress"`),
		),
		fx.Annotate(
			computers.NewAgifyComputer,
			fx.ParamTags(`name:"agifyAddress"`, ``),
			fx.As(new(messagehandler.AgeComputer)),
		),

		fx.Annotate(
			genderizeAddress,
			fx.ResultTags(`name:"genderizeAddress"`),
		),
		fx.Annotate(
			computers.NewSexComputer,
			fx.ParamTags(`name:"genderizeAddress"`, ``),
			fx.As(new(messagehandler.SexComputer)),
		),

		fx.Annotate(
			nationalityAddress,
			fx.ResultTags(`name:"nationalityAddress"`),
		),
		fx.Annotate(
			computers.NewNationalityComputer,
			fx.ParamTags(`name:"nationalityAddress"`, ``),
			fx.As(new(messagehandler.NationalityComputer)),
		),
	),
	fx.Supply(http.Client{}),
	fx.Decorate(decorateLogger),
	fx.Invoke(func(_ *service.EnrichService) {}),
)

type serviceParams struct {
	fx.In
	Runner   service.FioHandlingRunner
	Enricher *messagehandler.Enricher
	Logger   *zap.Logger
	Repo     service.Storage
}

func setupRunner(lc fx.Lifecycle, config *Config, logger *zap.Logger) (*handlingrunner.KafkaConsumerGroupRunner, error) {
	runner, err := handlingrunner.NewKafkaConsumerGroupRunner(config.DataBus.Brokers, config.DataBus.Topic, logger)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.StopHook(runner.Close))
	return runner, nil
}

func setupServiceLifecycle(lc fx.Lifecycle, params serviceParams) *service.EnrichService {
	s := service.NewEnrichService(params.Runner, params.Enricher, params.Logger, params.Repo)
	var cancelCtxFn context.CancelFunc
	done := make(chan struct{})
	lc.Append(fx.StartHook(func(ctx context.Context) {
		var childCtx context.Context
		childCtx, cancelCtxFn = context.WithCancel(ctx)
		go func() {
			defer close(done)
			s.Run(childCtx)
		}()
	}))
	lc.Append(fx.StopHook(func() {
		cancelCtxFn()
		<-done
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

func agifyAddress(config *Config) string {
	return config.AgifyAddress
}

func genderizeAddress(config *Config) string {
	return config.GenderizeAddress
}

func nationalityAddress(config *Config) string {
	return config.NationalityAddress
}
