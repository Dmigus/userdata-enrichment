package repository

import (
	"context"
	"enrichstorage/internal/service/enrichstorage/get"
	"enrichstorage/internal/service/enrichstorage/update"
	"enrichstorage/pkg/types"
	"errors"
	"time"

	"github.com/samber/lo"
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

func (r *Records) Get(ctx context.Context, req get.Request) ([]types.EnrichedRecord, error) {
	var dbRecords []Record
	db := r.db.WithContext(ctx)
	db = setWhereFromFilters(db, req.Filters)
	db = getPageFromPagination(db, req.Pagination)
	res := db.Find(&dbRecords)
	if res.Error != nil {
		return nil, res.Error
	}
	records := lo.Map(dbRecords, func(rec Record, _ int) types.EnrichedRecord {
		return dbRecordToEnrichedRecord(rec)
	})
	return records, nil
}

func (r *Records) DoesHaveRecordsBefore(ctx context.Context, filters get.Filters, fio types.FIO) (bool, error) {
	db := r.db.WithContext(ctx)
	db = setWhereFromFilters(db, filters)
	db = setBeforeCondition(db, fio)
	var count int64
	result := db.Model(&Record{}).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r *Records) DoesHaveRecordsAfter(ctx context.Context, filters get.Filters, fio types.FIO) (bool, error) {
	db := r.db.WithContext(ctx)
	db = setWhereFromFilters(db, filters)
	db = setAfterCondition(db, fio)
	var count int64
	result := db.Model(&Record{}).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
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
		if ageInterval.IsGtePresents() {
			db = db.Where("age >= ?", ageInterval.GetGTE())
		}
		if ageInterval.IsLtePresents() {
			db = db.Where("age <= ?", ageInterval.GetLTE())
		}
	}
	return db
}

func getPageFromPagination(db *gorm.DB, p get.Pagination) *gorm.DB {
	db = db.Limit(p.Limit)
	if lowLimit, ok := p.Before(); ok {
		db = db.Order("surname desc, name desc, patronymic desc")
		return setBeforeCondition(db, *lowLimit)
	} else if highLimit, ok := p.After(); ok {
		db = db.Order("surname asc, name asc, patronymic asc")
		return setAfterCondition(db, *highLimit)
	} else {
		db = db.Order("surname asc, name asc, patronymic asc")
	}
	return db
}

func setBeforeCondition(db *gorm.DB, fio types.FIO) *gorm.DB {
	return db.Where(`surname < ? OR surname = ? AND name < ? OR surname = ? AND name = ? AND patronymic < ?`,
		fio.Surname(), fio.Surname(), fio.Name(), fio.Surname(), fio.Name(), fio.Patronymic())
}

func setAfterCondition(db *gorm.DB, fio types.FIO) *gorm.DB {
	return db.Where(`surname > ? OR surname = ? AND name > ? OR surname = ? AND name = ? AND patronymic > ?`,
		fio.Surname(), fio.Surname(), fio.Name(), fio.Surname(), fio.Name(), fio.Patronymic())
}

func dbRecordToEnrichedRecord(rec Record) types.EnrichedRecord {
	fio, _ := types.NewFIO(rec.Fio.Name, rec.Fio.Surname, rec.Fio.Patronymic)
	return types.EnrichedRecord{
		Fio:         fio,
		Age:         rec.Age,
		Sex:         rec.Sex,
		Nationality: rec.Nationality,
	}
}

func keyFromModel(fio types.FIO) fioKey {
	return fioKey{Name: fio.Name(), Surname: fio.Surname(), Patronymic: fio.Patronymic()}
}
