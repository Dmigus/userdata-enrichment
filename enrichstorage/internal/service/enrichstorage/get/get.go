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
	before, after, err := g.compPaging(ctx, req, data)
	if err != nil {
		return Result{}, err
	}
	return Result{
		Records:  data,
		PrevPage: before,
		NextPage: after,
	}, nil
}

func (g *Getter) compPaging(ctx context.Context, req Request, data []types.EnrichedRecord) (*PrevPage, *NextPage, error) {
	var leftBorder, rightBorder *types.FIO
	if len(data) == 0 {
		leftBorder, _ = req.Pagination.After()
		rightBorder, _ = req.Pagination.Before()
	} else {
		minRecord := slices.MinFunc(data, func(a, b types.EnrichedRecord) int {
			return g.comparator.Cmp(a.Fio, b.Fio)
		})
		leftBorder = &minRecord.Fio
		maxRecord := slices.MaxFunc(data, func(a, b types.EnrichedRecord) int {
			return g.comparator.Cmp(a.Fio, b.Fio)
		})
		rightBorder = &maxRecord.Fio
	}
	return g.compPagingForWindow(ctx, req.Filters, leftBorder, rightBorder)
}

func (g *Getter) compPagingForWindow(ctx context.Context, filters Filters, leftBorder, rightBorder *types.FIO) (*PrevPage, *NextPage, error) {
	var prev *PrevPage
	if leftBorder != nil {
		have, err := g.repo.DoesHaveRecordsBefore(ctx, filters, *leftBorder)
		if err != nil {
			return nil, nil, err
		}
		if have {
			page := PrevPage(*leftBorder)
			prev = &page
		}
	}
	var after *NextPage
	if rightBorder != nil {
		have, err := g.repo.DoesHaveRecordsAfter(ctx, filters, *rightBorder)
		if err != nil {
			return nil, nil, err
		}
		if have {
			page := NextPage(*rightBorder)
			after = &page
		}
	}
	return prev, after, nil
}

func (g *Getter) IsFIOPresents(ctx context.Context, fio types.FIO) (bool, error) {
	return g.repo.IsFIOPresents(ctx, fio)
}
