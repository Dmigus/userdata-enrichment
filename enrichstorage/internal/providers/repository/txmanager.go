package repository

import (
	"context"
	"enrichstorage/internal/service/usecases/create"
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

func (m *TxManager) WithinTransaction(ctx context.Context, scenario func(context.Context, create.Records, create.Outbox) bool) error {
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

func clearExtraError(err error) error {
	if errors.Is(err, errScenarioChooseToNotCommit) {
		return nil
	}
	return err
}
