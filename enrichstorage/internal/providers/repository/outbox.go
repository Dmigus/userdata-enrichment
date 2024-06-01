package repository

import (
	"context"
	"encoding/json"
	"enrichstorage/pkg/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	FioOutbox struct {
		ID      int64 `gorm:"primaryKey"`
		Payload []byte
	}
	Outbox struct {
		db *gorm.DB
	}
)

func (o *Outbox) FIOComputeRequested(ctx context.Context, fio types.FIO) error {
	rec := FioOutbox{Payload: fioToBytes(fio)}
	result := o.db.WithContext(ctx).Select("payload").Create(&rec)
	return result.Error
}

func (o *Outbox) PullNextFIO(ctx context.Context, batchSize int) ([]types.FIO, error) {
	var result []types.FIO
	resDB := o.db.WithContext(ctx).
		Limit(batchSize).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Find(&result)
	if resDB.Error != nil {
		return nil, resDB.Error
	}
	return result, nil
}

func fioToBytes(fio types.FIO) []byte {
	res, _ := json.Marshal(fio)
	return res
}
