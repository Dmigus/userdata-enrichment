package app

import (
	"context"
	"enrichstorage/internal/controllers/grpc"
	"enrichstorage/internal/providers/repository"
	"enrichstorage/internal/service/enrichstorage/create"
	"enrichstorage/internal/service/enrichstorage/delete"
	"enrichstorage/internal/service/enrichstorage/update"
	"enrichstorage/internal/service/outboxsender"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

var Module = fx.Module("app",
	fx.Provide(
		outboxSender,
		create.NewCreator,
		delete.NewDeleter,
		update.NewUpdater,
		//get.NewGetter,
		fx.Annotate(
			batchSize,
			fx.ResultTags(`name:"batchSize"`),
		),
		fx.Annotate(
			batchInterval,
			fx.ResultTags(`name:"batchInterval"`)),
		getDB,
		fx.Annotate(
			repository.NewOutbox,
			fx.As(new(create.Outbox)),
			fx.As(new(outboxsender.Outbox)),
		),
		fx.Annotate(
			repository.NewRecords,
			fx.As(new(grpc.Service)),
			fx.As(new(update.Repository)),
			fx.As(new(create.Repository)),
			fx.As(new(delete.Repository)),
		),
		repository.NewTxManager,
	),
	fx.Invoke(func(_ *outboxsender.Service) {}),
)

type outboxParams struct {
	fx.In
	tx            outboxsender.TxManager
	broker        outboxsender.EventsPusher
	batchInterval int `name:"batchInterval"`
	batchSize     int `name:"batchSize"`
	logger        *zap.Logger
}

func batchSize(config *Config) int {
	return config.DataBus.BatchSize
}

func batchInterval(config *Config) int {
	return config.DataBus.BatchInterval
}

func outboxSender(lc fx.Lifecycle, params outboxParams) *outboxsender.Service {
	iterInterval := time.Duration(params.batchInterval) * time.Second
	sender := outboxsender.NewService(params.tx, params.broker, iterInterval, params.batchSize, params.logger)
	var lifecycleCtx context.Context
	var cancel func()
	done := make(chan struct{})
	lc.Append(fx.StartStopHook(func(ctx context.Context) {
		lifecycleCtx, cancel = context.WithCancel(ctx)
		go func() {
			defer close(done)
			sender.Run(lifecycleCtx)
		}()
	}, func() {
		cancel()
		<-done
	}))
	return sender
}

func getDB(lc fx.Lifecycle, config *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.Storage.GetPostgresDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	lc.Append(fx.StopHook(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}))
	return db, nil
}
