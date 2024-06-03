// Package outboxsender содержит реализацию работы transactional outbox
package outboxsender

import (
	"context"
	"enrichstorage/pkg/types"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type (
	// Outbox представляет из себя хранилище событий, из которого можно получать набор событий не более
	Outbox interface {
		PullNextFIO(ctx context.Context, batchSize int) ([]types.FIO, error)
	}
	eventsPusher interface {
		SendMessages(ctx context.Context, fios []types.FIO) error
	}
	txManager interface {
		WithinTransaction(context.Context, func(ctx context.Context, outbox Outbox) bool) error
	}
	// Service представляет сервис, осуществляющий регулярную выгрузку ивентов из хранилища событий в брокер
	Service struct {
		tx           txManager
		broker       eventsPusher
		iterInterval time.Duration
		batchSize    int
		logger       *zap.Logger
	}
)

// NewService создаёт новый Service.
func NewService(tx txManager, broker eventsPusher, iterInterval time.Duration, batchSize int, logger *zap.Logger) *Service {
	return &Service{
		tx:           tx,
		broker:       broker,
		iterInterval: iterInterval,
		batchSize:    batchSize,
		logger:       logger,
	}
}

// Run запускает обработку событий. Блокирующая. Остановиться, когда переданный контекст будет отменён
func (s *Service) Run(ctx context.Context) {
	t := time.NewTicker(s.iterInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			err := s.iteration(ctx)
			if err != nil {
				s.logger.Error("processing error", zap.Error(err))
			}
		}
	}
}

func (s *Service) iteration(serviceLiveCtx context.Context) error {
	iterCtx, cancel := context.WithCancel(serviceLiveCtx)
	defer cancel()
	var serviceErr error
	trErr := s.tx.WithinTransaction(iterCtx, func(ctx context.Context, outbox Outbox) bool {
		events, err := outbox.PullNextFIO(ctx, s.batchSize)
		if err != nil {
			serviceErr = fmt.Errorf("could not pull events from outbox: %w", err)
			return false
		}
		err = s.broker.SendMessages(ctx, events)
		if err != nil {
			serviceErr = fmt.Errorf("could not push events to broker: %w", err)
			return false
		}
		return true
	})
	if serviceErr != nil {
		return serviceErr
	}
	return trErr
}
