package repository

import (
	"context"
	"enrichstorage/internal/service/usecases/enter"
	"gorm.io/gorm"
)

type TxManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) WithinTransaction(ctx context.Context, scenario func(enter.Records, enter.Outbox) bool) {

}
