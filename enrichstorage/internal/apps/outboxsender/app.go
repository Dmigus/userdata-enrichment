package outboxsender

import (
	"context"
	"enrichstorage/internal/providers/kafka"
	"enrichstorage/internal/providers/repository"
	"enrichstorage/internal/service/outboxsender"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("outboxsender",
	fx.Provide(
		fx.Annotate(
			kafka.NewSender,
			fx.ParamTags(`name:"brokers"`, `name:"topic"`),
			fx.As(new(outboxsender.EventsPusher)),
		),
		namedRequestEventBusFields,
		outboxSender,
		fx.Annotate(
			repository.NewTxManager,
			fx.As(new(outboxsender.TxManager)),
		),
	),
	fx.Invoke(func(_ *outboxsender.Service) {}),
)

type requestEventBusConfig struct {
	fx.Out
	Brokers       []string `name:"brokers"`
	Topic         string   `name:"topic"`
	BatchSize     int      `name:"batchSize"`
	BatchInterval int      `name:"batchInterval"`
}

func namedRequestEventBusFields(config *Config) requestEventBusConfig {
	return requestEventBusConfig{
		Brokers:       config.RequestEventBus.Brokers,
		Topic:         config.RequestEventBus.Topic,
		BatchSize:     config.RequestEventBus.BatchSize,
		BatchInterval: config.RequestEventBus.BatchInterval,
	}
}

type outboxParams struct {
	fx.In
	Tx            outboxsender.TxManager
	Broker        outboxsender.EventsPusher
	BatchInterval int `name:"batchInterval"`
	BatchSize     int `name:"batchSize"`
	Logger        *zap.Logger
}

func outboxSender(lc fx.Lifecycle, params outboxParams) *outboxsender.Service {
	iterInterval := time.Duration(params.BatchInterval) * time.Second
	sender := outboxsender.NewService(params.Tx, params.Broker, iterInterval, params.BatchSize, params.Logger)
	var lifecycleCtx context.Context
	var cancel func()
	done := make(chan struct{})
	lc.Append(fx.StartStopHook(func() {
		lifecycleCtx, cancel = context.WithCancel(context.Background())
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
