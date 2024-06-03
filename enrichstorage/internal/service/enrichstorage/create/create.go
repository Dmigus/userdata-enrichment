package create

import (
	"context"
	"enrichstorage/pkg/types"
)

type (
	Records interface {
		Create(ctx context.Context, key types.FIO) error
	}
	Outbox interface {
		FIOComputeRequested(ctx context.Context, fio types.FIO) error
	}
	txManager interface {
		WithinTransaction(context.Context, func(Records, Outbox) bool) error
	}
	Enterer struct {
		txManager txManager
	}
)

func NewCreator(txManager txManager) *Enterer {
	return &Enterer{txManager: txManager}
}

func (e *Enterer) Create(ctx context.Context, fio types.FIO) error {
	var businessErr error
	txErr := e.txManager.WithinTransaction(ctx, func(rec Records, out Outbox) bool {
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
