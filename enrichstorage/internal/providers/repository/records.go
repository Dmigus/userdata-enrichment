package repository

import (
	"context"
	"enrichstorage/internal/service/enrichstorage/get"
	"enrichstorage/internal/service/enrichstorage/update"
	"enrichstorage/pkg/types"
	"errors"
	"github.com/samber/lo"
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
	dbRec := Record{Fio: keyFromModel(rec.Fio),
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
	dbRec := Record{Fio: keyFromModel(fio)}
	result := r.db.WithContext(ctx).Select(keyFields).Clauses(clause.OnConflict{DoNothing: true}).Create(&dbRec)
	return result.Error
}

func (r *Records) Delete(ctx context.Context, fio types.FIO) error {
	dbRec := Record{Fio: keyFromModel(fio)}
	result := r.db.WithContext(ctx).Delete(&dbRec)
	return result.Error
}

func (r *Records) Get(ctx context.Context, req get.Request) ([]get.Result, error) {
	var records []Record
	db := r.db.WithContext(ctx)
	db = setWhereFromFilters(db, req.Filters)
	db = getPageFromPagination(db, req.Pagination)
	res := db.Find(&records)
	if res.Error != nil {
		return nil, res.Error
	}
	results := lo.Map(records, func(rec Record, _ int) get.Result {
		return recordToGetResult(rec)
	})
	return results, nil
}

func setWhereFromFilters(db *gorm.DB, f get.Filters) *gorm.DB {
	if name, ok := f.NameFilter(); ok {
		db = db.Where("name = ?", name.Val)
	}
	if surname, ok := f.SurnameFilter(); ok {
		db = db.Where("surname = ?", surname.Val)
	}
	if patronymic, ok := f.PatronymicFilter(); ok {
		db = db.Where("patronymic = ?", patronymic.Val)
	}
	if sex, ok := f.SexFilter(); ok {
		db = db.Where("sex = ?", sex.Val)
	}
	if nat, ok := f.NatFilter(); ok {
		db = db.Where("nationality = ?", nat.Val)
	}
	if ageInterval, ok := f.AgeFilter(); ok {
		db = db.Where("age BETWEEN ? AND ?", ageInterval.LTE, ageInterval.GTE)
	}
	return db
}

func getPageFromPagination(db *gorm.DB, p get.Pagination) *gorm.DB {
	db = db.Limit(p.Limit)
	if lowLimit, ok := p.Before(); ok {
		db = db.Order("surname desc, name desc, patronymic desc")
		return db.Where("surname < ? OR surname = ? AND name < ? OR surname = ? AND name = ? AND patronymic < ?",
			lowLimit.Surname(), lowLimit.Surname(), lowLimit.Name(), lowLimit.Surname(), lowLimit.Name(), lowLimit.Patronymic())
	} else if highLimit, ok := p.After(); ok {
		db = db.Order("surname asc, name asc, patronymic asc")
		return db.Where("surname > ? OR surname = ? AND name > ? OR surname = ? AND name = ? AND patronymic > ?",
			highLimit.Surname(), highLimit.Surname(), highLimit.Name(), highLimit.Surname(), highLimit.Name(), highLimit.Patronymic())
	}
	return db
}

func recordToGetResult(rec Record) get.Result {
	fio, _ := types.NewFIO(rec.Fio.Name, rec.Fio.Surname, rec.Fio.Patronymic)
	return get.Result{
		Key:         fio,
		Age:         rec.Age,
		Sex:         rec.Sex,
		Nationality: rec.Nationality,
	}
}

func keyFromModel(fio types.FIO) fioKey {
	return fioKey{Name: fio.Name(), Surname: fio.Surname(), Patronymic: fio.Patronymic()}
}
