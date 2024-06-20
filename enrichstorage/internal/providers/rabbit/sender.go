package rabbit

import (
	"context"
	"enrichstorage/pkg/types"
	"errors"
	"fmt"

	lop "github.com/samber/lo/parallel"
	amqp "github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

var errNotConfirmed = fmt.Errorf("message publishing not confirmed")

type (
	Sender struct {
		queue  string
		logger *zap.Logger
		pub    *amqp.Publisher
	}
	RabbitCreds struct {
		Name     string
		Password string
	}
)

func NewSender(addr string, queue string, creds RabbitCreds, logger *zap.Logger) (*Sender, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s", creds.Name, creds.Password, addr)
	conn, err := amqp.NewConn(
		url,
	)
	if err != nil {
		return nil, err
	}
	pub, err := amqp.NewPublisher(conn, amqp.WithPublisherOptionsConfirm)
	if err != nil {
		return nil, err
	}
	return &Sender{queue: queue, pub: pub, logger: logger}, nil
}

func (s *Sender) SendMessages(ctx context.Context, fios []types.FIO) error {
	errs := lop.Map(fios, func(item types.FIO, _ int) error {
		return s.sendFIOToQueue(ctx, item)
	})
	if err := errors.Join(errs...); err != nil {
		return err
	}
	return nil
}

func (s *Sender) sendFIOToQueue(ctx context.Context, fio types.FIO) error {
	confirms, err := s.pub.PublishWithDeferredConfirmWithContext(
		ctx,
		fio.ToBytes(),
		[]string{s.queue},
		amqp.WithPublishOptionsContentType("text/plain"),
	)
	if err != nil {
		return err
	} else if len(confirms) == 0 || confirms[0] == nil {
		return errNotConfirmed
	}
	received, err := confirms[0].WaitContext(ctx)
	if err != nil {
		return err
	}
	if received {
		return nil
	}
	return errNotConfirmed
}

func (s *Sender) Close() {
	if s.pub != nil {
		s.pub.Close()
	}
}
