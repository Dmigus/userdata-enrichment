// Package kafka содержит Sender для записи в топик кафки
package kafka

import (
	"context"
	"encoding/json"
	"enrichstorage/pkg/types"
	"github.com/IBM/sarama"
	"github.com/samber/lo"
)

type messagesSender interface {
	SendMessages([]*sarama.ProducerMessage) error
}

// Sender это провайдер, который умеет отправлять сообщения в кафку
type Sender struct {
	topic    string
	producer messagesSender
}

// NewSender создайт новый Sender
func NewSender(brokers []string, topic string) (*Sender, error) {
	cfg := atLeastOnceConfig()
	syncProducer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	return &Sender{
		producer: syncProducer,
		topic:    topic,
	}, nil
}

// SendMessages синхронно отправляет сообщения events в брокер
func (p *Sender) SendMessages(_ context.Context, fios []types.FIO) error {
	saramaMessages := lo.Map(fios, func(fio types.FIO, _ int) *sarama.ProducerMessage {
		return p.modelMessageToSarama(fio)
	})
	return p.producer.SendMessages(saramaMessages)
}

func (p *Sender) modelMessageToSarama(fio types.FIO) *sarama.ProducerMessage {
	bytes, _ := json.Marshal(fio)
	return &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(bytes),
	}
}

func atLeastOnceConfig() *sarama.Config {
	c := sarama.NewConfig()
	// at least once
	c.Producer.RequiredAcks = sarama.WaitForAll
	c.Producer.Return.Successes = true
	c.Producer.Return.Errors = true
	return c
}
