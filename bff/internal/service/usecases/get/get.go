package get

import (
	"bff/pkg/types"
	"context"
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
	repo interface {
		Get(ctx context.Context, req Request) ([]Result, error)
	}
	Getter struct {
		repo repo
	}
)

func NewGetter(repo repo) *Getter {
	return &Getter{repo: repo}
}

func (g *Getter) Get(ctx context.Context, req Request) ([]Result, error) {
	return g.repo.Get(ctx, req)
}
