package kafka

import (
	"context"
	"enricher/internal/service"
	"enrichstorage/pkg/types"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"go.opentelemetry.io/otel"
)

type consumerGroupHandler struct {
	handler service.Handler
}

func newConsumerGroupHandler(handler service.Handler) *consumerGroupHandler {
	return &consumerGroupHandler{handler: handler}
}

// Setup Начинаем новую сессию, до ConsumeClaim
func (c *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup завершает сессию, после того, как все ConsumeClaim завершатся
func (c *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim обрабатывает до тех пор пока сессия не завершилась
func (c *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			c.processTracedMessage(session.Context(), message)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *consumerGroupHandler) processTracedMessage(sessionCtx context.Context, message *sarama.ConsumerMessage) {
	ctx := otel.GetTextMapPropagator().Extract(sessionCtx, otelsarama.NewConsumerMessageCarrier(message))
	fio := messageToFIO(message)
	c.handler.Handle(ctx, fio)
}

func messageToFIO(message *sarama.ConsumerMessage) types.FIO {
	fio, _ := types.FIOfromBytes(message.Value)
	return fio
}
