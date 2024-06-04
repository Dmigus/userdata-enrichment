package repository

import (
	"context"
	"enrichstorage/internal/service/enrichstorage/get"
	"enrichstorage/internal/service/enrichstorage/update"
	"enrichstorage/pkg/types"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	Record struct {
		fio         types.FIO `gorm:"embedded;primaryKey"`
		age         types.Age
		sex         types.Sex
		nationality types.Nationality
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	Records struct {
		db *gorm.DB
	}
)

func NewRecords(db *gorm.DB) *Records {
	return &Records{db: db}
}

func (Record) TableName() string {
	return "record"
}

func (r *Records) IsFIOPresents(ctx context.Context, fio types.FIO) (bool, error) {
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

func (r *Records) Update(ctx context.Context, rec update.Request) error {
	dbRec := Record{fio: rec.Fio, age: rec.NewAge, sex: rec.NewSex, nationality: rec.NewNat}
	var updateFields []string
	if rec.SexPresents {
		updateFields = append(updateFields, "sex")
	}
	if rec.AgePresents {
		updateFields = append(updateFields, "age")
	}
	if rec.NationalityPresents {
		updateFields = append(updateFields, "nationality")
	}
	result := r.db.WithContext(ctx).Select(updateFields).Updates(&dbRec)
	return result.Error
}

func (r *Records) Create(ctx context.Context, fio types.FIO) error {
	dbRec := Record{fio: fio}
	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&dbRec)
	return result.Error
}

func (r *Records) Delete(ctx context.Context, fio types.FIO) error {
	dbRec := Record{fio: fio}
	result := r.db.WithContext(ctx).Delete(&dbRec)
	return result.Error
}

func (r *Records) Get(ctx context.Context, req get.Request) ([]get.Result, error) {
	return nil, nil
}
