package create

import (
	"context"
	"enrichstorage/pkg/types"
)

type (
	Repository interface {
		Create(ctx context.Context, key types.FIO) error
	}
	Outbox interface {
		FIOComputeRequested(ctx context.Context, fio types.FIO) error
	}
	TxManager interface {
		WithinTransactionRecordsOutbox(context.Context, func(context.Context, Repository, Outbox) bool) error
	}
	Enterer struct {
		txManager TxManager
	}
)

func NewCreator(txManager TxManager) *Enterer {
	return &Enterer{txManager: txManager}
}

func (e *Enterer) Create(ctx context.Context, fio types.FIO) error {
	var businessErr error
	txErr := e.txManager.WithinTransactionRecordsOutbox(ctx, func(ctx context.Context, rec Repository, out Outbox) bool {
		err := rec.Create(ctx, fio)
		if err != nil {
			businessErr = err
			return false
		}
		err = out.FIOComputeRequested(ctx, fio)
		if err != nil {
			businessErr = err
			return false
		}
		return true
	})
	if businessErr != nil {
		return businessErr
	}
	return txErr
}
