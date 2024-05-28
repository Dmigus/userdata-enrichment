package service

import (
	"bff/pkg/types"
	"context"
	"go.uber.org/zap"
)

type (
	funcHandler func(context.Context, types.FIO)
	Handler     interface {
		Handle(ctx context.Context, fio types.FIO)
	}
	FioHandlingRunner interface {
		Run(context.Context, Handler) error
	}
	Repository interface {
		IsFIOPresents(ctx context.Context, fio types.FIO) (bool, error)
		Store(ctx context.Context, rec types.EnrichedRecord) error
	}
	Enricher interface {
		Enrich(context.Context, types.FIO) (types.EnrichedRecord, error)
	}
	EnrichService struct {
		enricher Enricher
		runner   FioHandlingRunner
		logger   *zap.Logger
		repo     Repository
	}
)

func NewEnrichService(runner FioHandlingRunner, enricher Enricher, logger *zap.Logger, repo Repository) *EnrichService {
	return &EnrichService{runner: runner, enricher: enricher, logger: logger, repo: repo}
}

func (en *EnrichService) Run(ctx context.Context) error {
	var handleFIOScenario funcHandler = func(ctx context.Context, fio types.FIO) {
		present, err := en.repo.IsFIOPresents(ctx, fio)
		if err != nil {
			en.handleErr(ctx, fio, "err checking fio presence", err)
			return
		}
		if !present {
			en.logger.Info("fio is not present in repository", types.FioToZaFields(fio)...)
			return
		}
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
