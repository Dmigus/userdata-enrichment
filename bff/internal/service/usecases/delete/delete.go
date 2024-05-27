package delete

import (
	"bff/pkg/types"
	"context"
)

type (
	repo interface {
		Delete(ctx context.Context, key types.FIO) error
	}
	Deleter struct {
		repo repo
	}
)

func NewDeleter(repo repo) *Deleter {
	return &Deleter{repo: repo}
}

func (e *Deleter) Delete(ctx context.Context, fio types.FIO) error {
	return e.repo.Delete(ctx, fio)
}
