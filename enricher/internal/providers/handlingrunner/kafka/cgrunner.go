// Package handlingrunner содержит функциональность, которая позволяет запустить обработку событий, получаемых из некоторого источника
package kafka

import (
	"context"
	"enricher/internal/service"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const groupName = "enricher-group"

// KafkaConsumerGroupRunner это структура, которая умеет запускать обработку событий, получаемых из кафки
type KafkaConsumerGroupRunner struct {
	topic  string
	logger *zap.Logger
	cg     sarama.ConsumerGroup
}

// NewKafkaConsumerGroupRunner возращает новый KafkaConsumerGroupRunner, сконфигурированный на брокеры brokers и топик topic
func NewKafkaConsumerGroupRunner(brokers []string, topic string, logger *zap.Logger) (*KafkaConsumerGroupRunner, error) {
	cg, err := sarama.NewConsumerGroup(brokers, groupName, getConfig())
	if err != nil {
		return nil, err
	}
	return &KafkaConsumerGroupRunner{
		cg:     cg,
		topic:  topic,
		logger: logger,
	}, nil
}

// Run обрабатывает поступающие сообщения переданным хандлером в рамках группы. Блокирующий.
func (k *KafkaConsumerGroupRunner) Run(ctx context.Context, handler service.Handler) {
	saramaHandler := newConsumerGroupHandler(handler)
	for {
		err := k.cg.Consume(ctx, []string{k.topic}, saramaHandler)
		if err != nil {
			k.logger.Error("error in consumer group session", zap.Error(err))
		}
		if ctx.Err() != nil {
			return
		}
	}
}

func (k *KafkaConsumerGroupRunner) Close() error {
	return k.cg.Close()
}

func getConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	return config
}
