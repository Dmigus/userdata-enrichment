package repository

import (
	"bff/pkg/types"
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type (
	Record struct {
		fio         types.FIO `gorm:"embedded;primaryKey;->;<-:false"`
		age         types.Age
		sex         types.Sex
		nationality types.Nationality
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
	Repository struct {
		db gorm.DB
	}
)

func (Record) TableName() string {
	return "Record"
}

func (r *Repository) IsFIOPresents(ctx context.Context, fio types.FIO) (bool, error) {
	rec := Record{fio: fio}
	result := r.db.WithContext(ctx).First(&rec)
	if result.Error == nil {
		return true, nil
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, result.Error
}

func (r *Repository) Store(ctx context.Context, rec types.EnrichedRecord) error {
	dbRec := Record{fio: rec.Key, age: rec.Age, sex: rec.Sex, nationality: rec.Nationality}
	result := r.db.WithContext(ctx).Save(&dbRec)
	return result.Error
}
