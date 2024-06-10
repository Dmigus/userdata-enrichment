package get

import (
	"context"
	"enrichstorage/pkg/types"
	"slices"
)

type (
	Request struct {
		Filters    Filters
		Pagination Pagination
	}
	PrevPage types.FIO
	NextPage types.FIO
	Result   struct {
		Records  []types.EnrichedRecord
		PrevPage *PrevPage
		NextPage *NextPage
	}
	KeysComparator interface {
		Cmp(a, b types.FIO) int
	}
	Repository interface {
		Get(context.Context, Request) ([]types.EnrichedRecord, error)
		DoesHaveRecordsBefore(context.Context, Filters, types.FIO) (bool, error)
		DoesHaveRecordsAfter(context.Context, Filters, types.FIO) (bool, error)
		IsFIOPresents(context.Context, types.FIO) (bool, error)
	}
	Getter struct {
		repo       Repository
		comparator KeysComparator
	}
)

func NewGetter(repo Repository, comparator KeysComparator) *Getter {
	return &Getter{repo: repo, comparator: comparator}
}

func (g *Getter) GetWithPaging(ctx context.Context, req Request) (Result, error) {
	data, err := g.repo.Get(ctx, req)
	if err != nil {
		return Result{}, err
	}
	var before *PrevPage
	var after *NextPage
	if len(data) == 0 {
		before, after, err = g.compPagingForEmpty(ctx, req)
	} else {
		before, after, err = g.compPagingForExisting(ctx, req, data)
	}
	if err != nil {
		return Result{}, err
	}
	return Result{
		Records:  data,
		PrevPage: before,
		NextPage: after,
	}, nil
}

func (g *Getter) compPagingForEmpty(ctx context.Context, req Request) (before *PrevPage, after *NextPage, err error) {
	if bef, ok := req.Pagination.Before(); ok {
		have, err := g.repo.DoesHaveRecordsAfter(ctx, req.Filters, *bef)
		if err != nil {
			return nil, nil, err
		}
		if have {
			page := NextPage(*bef)
			return nil, &page, nil
		}
	} else if aft, ok := req.Pagination.After(); ok {
		have, err := g.repo.DoesHaveRecordsBefore(ctx, req.Filters, *aft)
		if err != nil {
			return nil, nil, err
		}
		if have {
			page := PrevPage(*aft)
			return &page, nil, nil
		}
	}
	return nil, nil, nil
}

func (g *Getter) compPagingForExisting(ctx context.Context, req Request, res []types.EnrichedRecord) (before *PrevPage, after *NextPage, err error) {
	maxRecord := slices.MaxFunc(res, func(a, b types.EnrichedRecord) int {
		return g.comparator.Cmp(a.Fio, b.Fio)
	})
	have, err := g.repo.DoesHaveRecordsAfter(ctx, req.Filters, maxRecord.Fio)
	if err != nil {
		return nil, nil, err
	}
	if have {
		page := NextPage(maxRecord.Fio)
		after = &page
	}
	minRecord := slices.MinFunc(res, func(a, b types.EnrichedRecord) int {
		return g.comparator.Cmp(a.Fio, b.Fio)
	})
	have, err = g.repo.DoesHaveRecordsBefore(ctx, req.Filters, minRecord.Fio)
	if err != nil {
		return nil, nil, err
	}
	if have {
		page := PrevPage(minRecord.Fio)
		before = &page
	}
	return before, after, nil
}

func (g *Getter) IsFIOPresents(ctx context.Context, fio types.FIO) (bool, error) {
	return g.repo.IsFIOPresents(ctx, fio)
}
