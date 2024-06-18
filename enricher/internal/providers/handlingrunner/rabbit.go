package handlingrunner

import (
	"context"
	"enricher/internal/service"
	"enrichstorage/pkg/types"
	"errors"
	"fmt"
	"slices"

	advancedErrs "github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/samber/lo"
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
	msgs, closeRes, err := r.getMessagesChan(ctx)
	if err != nil {
		r.logger.Error("error initializing consuming", zap.Error(err))
		return
	}
	defer closeRes.close()
	for msg := range msgs {
		fio, err := types.FIOfromBytes(msg.Body)
		if err != nil {
			_ = msg.Nack(true, false)
		} else {
			handler.Handle(ctx, fio)
		}
	}
}

func (r *RabbitRunner) getMessagesChan(ctx context.Context) (<-chan amqp.Delivery, *resToClose, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s", r.creds.Name, r.creds.Password, r.addr)
	conn, err := amqp.Dial(url)
	closeRes := &resToClose{}
	if err != nil {
		return nil, nil, advancedErrs.Wrap(err, "error connecting to rabbit")
	}
	closeRes.append(conn.Close)
	ch, err := conn.Channel()
	if err != nil {
		_ = closeRes.close()
		return nil, nil, advancedErrs.Wrap(err, "error getting channel to rabbit")
	}
	closeRes.append(ch.Close)
	msgs, err := ch.ConsumeWithContext(ctx, r.queue, "",
		false,
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		_ = closeRes.close()
		return nil, nil, advancedErrs.Wrap(err, "error getting messages channel")
	}
	return msgs, closeRes, err
}

type (
	resToClose struct {
		closeFns []func() error
	}
)

func (rc *resToClose) append(f func() error) {
	rc.closeFns = append(rc.closeFns, f)
}

func (rc *resToClose) close() error {
	slices.Reverse(rc.closeFns)
	errs := lo.Map(rc.closeFns, func(item func() error, _ int) error {
		return item()
	})
	return errors.Join(errs...)
}
