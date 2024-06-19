package rabbit

import (
	"context"
	"enricher/internal/service"
	"enrichstorage/pkg/types"
	"fmt"
	advancedErrs "github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
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
	msgs, closeFunc, err := r.getMessagesChan(ctx)
	if err != nil {
		r.logger.Error("error initializing consuming", zap.Error(err))
		return
	}
	defer closeFunc()
	for msg := range msgs {
		fio, err := types.FIOfromBytes(msg.Body)
		if err != nil {
			_ = msg.Nack(false, false)
		} else {
			handler.Handle(ctx, fio)
			_ = msg.Ack(false)
		}
	}
}

func (r *RabbitRunner) getMessagesChan(ctx context.Context) (<-chan amqp.Delivery, func() error, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s", r.creds.Name, r.creds.Password, r.addr)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, advancedErrs.Wrap(err, "error connecting to rabbit")
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, nil, advancedErrs.Wrap(err, "error getting channel to rabbit")
	}
	msgs, err := ch.ConsumeWithContext(ctx, r.queue, "",
		false,
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		_ = conn.Close()
		return nil, nil, advancedErrs.Wrap(err, "error getting messages channel")
	}
	return msgs, conn.Close, err
}
