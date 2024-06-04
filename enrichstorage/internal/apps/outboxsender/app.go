package outboxsender

import (
	"context"
	"enrichstorage/internal/service/outboxsender"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("outboxsender",
	fx.Provide(
		outboxSender,
		fx.Annotate(
			batchSize,
			fx.ResultTags(`name:"batchSize"`),
		),
		fx.Annotate(
			batchInterval,
			fx.ResultTags(`name:"batchInterval"`)),
	),
	fx.Invoke(func(_ *outboxsender.Service) {}),
)

func batchSize(config *Config) int {
	return config.DataBus.BatchSize
}

func batchInterval(config *Config) int {
	return config.DataBus.BatchInterval
}

type outboxParams struct {
	fx.In
	tx            outboxsender.TxManager
	broker        outboxsender.EventsPusher
	batchInterval int `name:"batchInterval"`
	batchSize     int `name:"batchSize"`
	logger        *zap.Logger
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
