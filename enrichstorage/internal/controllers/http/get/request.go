package get

import (
	"enrichstorage/internal/service/enrichstorage/get"
	"enrichstorage/pkg/types"
)

const defaultLimit = 2

type request struct {
	Surname     *string `form:"surname"`
	Name        *string `form:"name"`
	Patronymic  *string `form:"patronymic"`
	Sex         *string `form:"sex"`
	Nationality *string `form:"nationality"`
	AgeGte      *int    `form:"age[gte]"`
	AgeLte      *int    `form:"age[lte]"`
	Limit       *int    `form:"limit"`
	After       *string `form:"after"`
	Before      *string `form:"before"`
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
	limit := defaultLimit
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
