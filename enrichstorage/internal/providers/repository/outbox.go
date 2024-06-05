package repository

import (
	"context"
	"enrichstorage/pkg/types"

	"github.com/samber/lo"
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

func (FioOutbox) TableName() string {
	return "fio_outbox"
}

func NewOutbox(db *gorm.DB) *Outbox {
	return &Outbox{db: db}
}

func (o *Outbox) FIOComputeRequested(ctx context.Context, fio types.FIO) error {
	rec := FioOutbox{Payload: fio.ToBytes()}
	result := o.db.WithContext(ctx).Select("payload").Create(&rec)
	return result.Error
}

func (o *Outbox) PullNextFIO(ctx context.Context, batchSize int) ([]types.FIO, error) {
	idsToDelete := []int64{}
	resDB := o.db.WithContext(ctx).
		Model(&FioOutbox{}).
		Select("id").
		Limit(batchSize).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Order("id").
		Find(&idsToDelete)
	if resDB.Error != nil {
		return nil, resDB.Error
	}
	var resultBytes []FioOutbox
	resDB = o.db.WithContext(ctx).
		Where("id IN ?", idsToDelete).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "payload"}}}).
		Delete(&resultBytes)
	if resDB.Error != nil {
		return nil, resDB.Error
	}
	fios := lo.Map(resultBytes, func(bytes FioOutbox, _ int) types.FIO {
		fio, _ := types.FIOfromBytes(bytes.Payload)
		return fio
	})
	return fios, nil
}
