package rabbit

import (
	"context"
	"enrichstorage/pkg/types"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	lop "github.com/samber/lo/parallel"
	"go.uber.org/zap"
)

type (
	Sender struct {
		queue  string
		logger *zap.Logger
		conn   *amqp.Connection
	}
	RabbitCreds struct {
		Name     string
		Password string
	}
)

func NewSender(addr string, queue string, creds RabbitCreds, logger *zap.Logger) (*Sender, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s", creds.Name, creds.Password, addr)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Sender{queue: queue, conn: conn, logger: logger}, nil
}

func (s *Sender) SendMessages(ctx context.Context, fios []types.FIO) error {
	ch, err := s.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	published := make(chan amqp.Confirmation, len(fios)+1)
	confirmed := ch.NotifyPublish(published)
	if err = ch.Confirm(false); err != nil {
		return err
	}
	errs := lop.Map(fios, func(item types.FIO, _ int) error {
		return s.sendFIOToQueue(ctx, item, ch)
	})
	_ = ch.Close()
	if err = errors.Join(errs...); err != nil {
		return err
	}
	for confirmation := range confirmed {
		if !checkConfirmation(confirmation) {
			return fmt.Errorf("at least one of fios was not published")
		}
	}
	return nil
}

func (s *Sender) sendFIOToQueue(ctx context.Context, fio types.FIO, ch *amqp.Channel) error {
	return ch.PublishWithContext(ctx,
		"",      // exchange
		s.queue, // routing key
		true,    // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        fio.ToBytes(),
		})
}

func checkConfirmation(c amqp.Confirmation) bool {
	return c.Ack
}

func (s *Sender) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
