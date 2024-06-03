package get

import (
	"context"
	"enrichstorage/pkg/types"
)

type (
	Request struct {
	}
	Result struct {
		Key         types.FIO
		Age         types.Age
		Sex         types.Sex
		Nationality types.Nationality
	}
	Repository interface {
		Get(ctx context.Context, req Request) ([]Result, error)
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
