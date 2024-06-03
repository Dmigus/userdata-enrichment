package repository

import (
	"context"
	"enrichstorage/internal/service/enrichstorage/create"
	"errors"

	"gorm.io/gorm"
)

var errScenarioChooseToNotCommit = errors.New("scenario choose to not commit")

type TxManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) WithinTransactionRecordsOutbox(ctx context.Context, scenario func(context.Context, create.Repository, create.Outbox) bool) error {
	err := m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx := tx.Statement.Context
		records := NewRecords(tx)
		outbox := NewOutbox(tx)
		if scenario(ctx, records, outbox) {
			return nil
		} else {
			// необходимо создать ошибку, чтобы был rollback
			return errScenarioChooseToNotCommit
		}
	})
	err = clearExtraError(err)
	return err
}

func (m *TxManager) WithinTransactionOutbox(ctx context.Context, scenario func(context.Context, create.Outbox) bool) error {
	err := m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx := tx.Statement.Context
		outbox := NewOutbox(tx)
		if scenario(ctx, outbox) {
			return nil
		} else {
			// необходимо создать ошибку, чтобы был rollback
			return errScenarioChooseToNotCommit
		}
	})
	err = clearExtraError(err)
	return err
}

func clearExtraError(err error) error {
	if errors.Is(err, errScenarioChooseToNotCommit) {
		return nil
	}
	return err
}
