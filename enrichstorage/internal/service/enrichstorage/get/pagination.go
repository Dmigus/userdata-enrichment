package get

import "enrichstorage/pkg/types"

// Pagination это объект, который хранит достаточно информации для формирования одной страницы данных
type Pagination struct {
	Limit         int
	before, after *types.FIO
}

func NewPagination(limit int, before *types.FIO, after *types.FIO) *Pagination {
	return &Pagination{Limit: limit, before: before, after: after}
}

func (p *Pagination) Before() (*types.FIO, bool) {
	return p.before, p.before != nil
}

func (p *Pagination) After() (*types.FIO, bool) {
	return p.after, p.after != nil
}
