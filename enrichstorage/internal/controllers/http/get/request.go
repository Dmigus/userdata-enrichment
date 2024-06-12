package get

import (
	"enrichstorage/internal/service/enrichstorage/get"
	"enrichstorage/pkg/types"
)

const defaultPageSize = 2

type request struct {
	Surname     *string `form:"surname"  validate:"optional"`
	Name        *string `form:"name" validate:"optional"`
	Patronymic  *string `form:"patronymic" validate:"optional"`
	Sex         *string `form:"sex" validate:"optional"`
	Nationality *string `form:"nationality" validate:"optional"`
	AgeGte      *int    `form:"age[gte]" validate:"optional"`
	AgeLte      *int    `form:"age[lte]" validate:"optional"`
	Limit       *int    `form:"limit" validate:"optional"`
	After       *string `form:"after" validate:"optional"`
	Before      *string `form:"before" validate:"optional"`
}

func (r *request) ToUsecaseRequest() (*get.Request, error) {
	result := &get.Request{}
	result.Filters = r.filters()
	pagination, err := r.pagination()
	if err != nil {
		return nil, err
	}
	result.Pagination = *pagination
	return result, nil
}

func (r *request) filters() get.Filters {
	filters := get.NewFilters()
	if r.Surname != nil {
		filters.SetSurnameFilter(&get.SurnameFilter{Val: *r.Surname})
	}
	if r.Name != nil {
		filters.SetNameFilter(&get.NameFilter{Val: *r.Name})
	}
	if r.Patronymic != nil {
		filters.SetPatronymicFilter(&get.PatronymicFilter{Val: *r.Patronymic})
	}
	if r.Sex != nil {
		filters.SetSexFilter(&get.SexFilter{Val: *r.Sex})
	}
	if r.Nationality != nil {
		filters.SetNationalityFilter(&get.NationalityFilter{Val: *r.Sex})
	}
	if r.AgeGte != nil || r.AgeLte != nil {
		filter := &get.AgeFilter{}
		if r.AgeGte != nil {
			filter.GTE = r.AgeGte
		}
		if r.AgeLte != nil {
			filter.LTE = r.AgeLte
		}
		filters.SetAgeFilter(filter)
	}
	return *filters
}

func (r *request) pagination() (*get.Pagination, error) {
	limit := defaultPageSize
	if r.Limit != nil {
		limit = *r.Limit
	}
	var afterFIO, beforeFIO *types.FIO
	if r.After != nil {
		fio, err := unmarshalFIO(*r.After)
		if err != nil {
			return nil, err
		}
		afterFIO = fio
	}
	if r.Before != nil {
		fio, err := unmarshalFIO(*r.Before)
		if err != nil {
			return nil, err
		}
		beforeFIO = fio
	}
	return get.NewPagination(limit, beforeFIO, afterFIO), nil
}
