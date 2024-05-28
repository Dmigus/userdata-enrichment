package service

import (
	"bff/pkg/types"
	"context"
	"enricher/internal/providers/messagehandler"
	"go.uber.org/zap"
)

type (
	funcHandler func(context.Context, types.FIO)
	Handler     interface {
		Handle(ctx context.Context, fio types.FIO)
	}
	fioHandlingRunner interface {
		Run(context.Context, Handler) error
	}
	repository interface {
		Store(ctx context.Context, rec types.EnrichedRecord) error
	}
	EnrichService struct {
		enricher messagehandler.Enricher
		runner   fioHandlingRunner
		logger   *zap.Logger
		repo     repository
	}
)

func NewEnrichService(runner fioHandlingRunner, enricher messagehandler.Enricher, logger *zap.Logger, repo repository) *EnrichService {
	return &EnrichService{runner: runner, enricher: enricher, logger: logger, repo: repo}
}

func (en *EnrichService) Run(ctx context.Context) error {
	var handleFIOScenario funcHandler = func(ctx context.Context, fio types.FIO) {
		enriched, err := en.enricher.Enrich(ctx, fio)
		if err != nil {
			en.handleErr(ctx, fio, "err enriching fio", err)
			return
		}
		err = en.repo.Store(ctx, enriched)
		en.handleErr(ctx, fio, "err storing to repository", err)
	}
	return en.runner.Run(ctx, handleFIOScenario)
}

func (en *EnrichService) handleErr(ctx context.Context, fio types.FIO, msg string, err error) {
	if err != nil {
		fields := append(types.FioToZaFields(fio), zap.Error(err))
		en.logger.Error(msg, fields...)
	}
}

func (fh funcHandler) Handle(ctx context.Context, fio types.FIO) {
	fh(ctx, fio)
}
