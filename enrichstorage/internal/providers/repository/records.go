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

var keyFields = []string{"surname", "name", "patronymic"}

type (
	fioKey struct {
		Name       string `gorm:"primaryKey"`
		Surname    string `gorm:"primaryKey"`
		Patronymic string `gorm:"primaryKey"`
	}
	Record struct {
		Fio         fioKey `gorm:"embedded"`
		Age         types.Age
		Sex         types.Sex
		Nationality types.Nationality
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
	rec := Record{Fio: fioKey{Name: fio.Name(), Surname: fio.Surname(), Patronymic: fio.Patronymic()}}
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
	dbRec := Record{Fio: fioKey{
		Name:       rec.Fio.Name(),
		Surname:    rec.Fio.Surname(),
		Patronymic: rec.Fio.Patronymic()},
		Age:         rec.NewAge,
		Sex:         rec.NewSex,
		Nationality: rec.NewNat,
	}
	updateFields := make(map[string]any)
	if rec.SexPresents {
		updateFields["sex"] = rec.NewSex
	}
	if rec.AgePresents {
		updateFields["age"] = rec.NewAge
	}
	if rec.NationalityPresents {
		updateFields["nationality"] = rec.NewNat
	}
	result := r.db.WithContext(ctx).Model(&dbRec).Updates(updateFields)
	return result.Error
}

func (r *Records) Create(ctx context.Context, fio types.FIO) error {
	dbRec := Record{Fio: fioKey{Name: fio.Name(), Surname: fio.Surname(), Patronymic: fio.Patronymic()}}
	result := r.db.WithContext(ctx).Select(keyFields).Clauses(clause.OnConflict{DoNothing: true}).Create(&dbRec)
	return result.Error
}

func (r *Records) Delete(ctx context.Context, fio types.FIO) error {
	dbRec := Record{Fio: fioKey{Name: fio.Name(), Surname: fio.Surname(), Patronymic: fio.Patronymic()}}
	result := r.db.WithContext(ctx).Delete(&dbRec)
	return result.Error
}

func (r *Records) Get(ctx context.Context, req get.Request) ([]get.Result, error) {
	return nil, nil
}
