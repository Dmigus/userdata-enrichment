package rabbit

import (
	"context"
	"enricher/internal/service"
	"enrichstorage/pkg/types"
	"fmt"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
	"sync"
)

type (
	RabbitRunner struct {
		addr   string
		queue  string
		creds  RabbitCreds
		logger *zap.Logger
	}
	RabbitCreds struct {
		Name     string
		Password string
	}
)

func NewRabbitRunner(addr string, queue string, creds RabbitCreds, logger *zap.Logger) (*RabbitRunner, error) {
	return &RabbitRunner{addr: addr, queue: queue, creds: creds, logger: logger}, nil
}

func (r *RabbitRunner) Run(ctx context.Context, handler service.Handler) {
	cons, err := newRabbitConsumer(r.addr, r.queue, r.creds, r.logger)
	if err != nil {
		r.logger.Error("error creating consumer client", zap.Error(err))
		return
	}
	closeOnce := sync.Once{}
	defer closeOnce.Do(cons.Close)
	go func() {
		<-ctx.Done()
		closeOnce.Do(cons.Close)
	}()
	err = cons.Run(func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		fio, err := types.FIOfromBytes(d.Body)
		if err != nil {
			return rabbitmq.NackDiscard
		}
		handler.Handle(ctx, fio)
		return rabbitmq.Ack
	})
	if err != nil {
		r.logger.Error("error during consuming", zap.Error(err))
	}
}

func newRabbitConsumer(addr string, queue string, creds RabbitCreds, logger *zap.Logger) (*rabbitmq.Consumer, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s", creds.Name, creds.Password, addr)
	conn, err := rabbitmq.NewConn(
		url,
	)
	if err != nil {
		return nil, err
	}
	consumer, err := rabbitmq.NewConsumer(
		conn,
		queue,
		rabbitmq.WithConsumerOptionsQueueNoDeclare,
		rabbitmq.WithConsumerOptionsLogger(logger.Sugar()),
	)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return consumer, nil
}
