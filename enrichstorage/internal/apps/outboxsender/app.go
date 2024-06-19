package outboxsender

import (
	"context"
	"enrichstorage/internal/providers/rabbit"
	"enrichstorage/internal/providers/repository"
	"enrichstorage/internal/service/outboxsender"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("outboxsender",
	fx.Provide(
		newSender,
		outboxSender,
		fx.Annotate(
			repository.NewTxManager,
			fx.As(new(outboxsender.TxManager)),
		),
	),
	fx.Invoke(func(_ *outboxsender.Service) {}),
)

func newSender(lc fx.Lifecycle, config *Config, logger *zap.Logger) (outboxsender.EventsPusher, error) {
	sender, err := rabbit.NewSender(
		config.RequestEventBus.Brokers[0],
		config.RequestEventBus.Topic,
		rabbit.RabbitCreds{Name: config.RequestEventBus.Username, Password: config.RequestEventBus.Password},
		logger,
	)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.StopHook(func() { sender.Close() }))
	return sender, nil
}

type outboxParams struct {
	fx.In
	Tx     outboxsender.TxManager
	Broker outboxsender.EventsPusher
	Config *Config
	Logger *zap.Logger
}

func outboxSender(lc fx.Lifecycle, params outboxParams) *outboxsender.Service {
	iterInterval := time.Duration(params.Config.RequestEventBus.BatchInterval) * time.Second
	sender := outboxsender.NewService(params.Tx, params.Broker, iterInterval, params.Config.RequestEventBus.BatchSize, params.Logger)
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
