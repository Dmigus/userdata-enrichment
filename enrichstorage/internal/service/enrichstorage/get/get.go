package get

import (
	"context"
	"enrichstorage/pkg/types"
)

type (
	Request struct {
		Filters    Filters
		Pagination Pagination
	}
	Result struct {
		Key         types.FIO
		Age         types.Age
		Sex         types.Sex
		Nationality types.Nationality
	}
	Repository interface {
		Get(ctx context.Context, req Request) ([]Result, error)
		IsFIOPresents(ctx context.Context, fio types.FIO) (bool, error)
	}
	Getter struct {
		repo Repository
	}
)

func NewGetter(repo Repository) *Getter {
	return &Getter{repo: repo}
}

func (g *Getter) Get(ctx context.Context, req Request) ([]Result, error) {
	return g.repo.Get(ctx, req)
}

func (g *Getter) IsFIOPresents(ctx context.Context, fio types.FIO) (bool, error) {
	return g.repo.IsFIOPresents(ctx, fio)
}
