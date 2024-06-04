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
		namedDatabusFields,
		outboxSender,
		fx.Annotate(
			repository.NewTxManager,
			fx.As(new(outboxsender.TxManager)),
		),
		//fx.Annotate(
		//	batchSize,
		//	fx.ResultTags(`name:"batchSize"`),
		//),
		//fx.Annotate(
		//	batchInterval,
		//	fx.ResultTags(`name:"batchInterval"`)),
	),
	fx.Invoke(func(_ *outboxsender.Service) {}),
)

type databusResult struct {
	fx.Out
	Brokers       []string `name:"brokers"`
	Topic         string   `name:"topic"`
	BatchSize     int      `name:"batchSize"`
	BatchInterval int      `name:"batchInterval"`
}

func namedDatabusFields(config *Config) databusResult {
	return databusResult{
		Brokers:       config.DataBus.Brokers,
		Topic:         config.DataBus.Topic,
		BatchSize:     config.DataBus.BatchSize,
		BatchInterval: config.DataBus.BatchInterval,
	}
}

//func batchSize(config *Config) int {
//	return config.DataBus.BatchSize
//}
//
//func batchInterval(config *Config) int {
//	return config.DataBus.BatchInterval
//}

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