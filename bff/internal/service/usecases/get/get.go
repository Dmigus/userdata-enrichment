package get

import (
	"bff/internal/service"
	"context"
)

type (
	Request struct {
	}
	Result struct {
		Key         service.FIO
		Age         service.Age
		Sex         service.Sex
		Nationality service.Nationality
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
