package service

import (
	"bff/pkg/types"
	"context"
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
	}
	EnrichService struct {
		runner fioHandlingRunner
	}
)

func NewEnrichService(runner fioHandlingRunner) *EnrichService {
	return &EnrichService{runner: runner}
}

func (en *EnrichService) Run(ctx context.Context) error {
	var handleFIOScenario funcHandler = func(ctx context.Context, fio types.FIO) {

	}
	return en.runner.Run(ctx, handleFIOScenario)
}

func (fh funcHandler) Handle(ctx context.Context, fio types.FIO) {
	fh(ctx, fio)
}
