package delete

import (
	"context"
	"enrichstorage/pkg/types"
)

type (
	Repository interface {
		Delete(ctx context.Context, key types.FIO) error
	}
	Deleter struct {
		repo Repository
	}
)

func NewDeleter(repo Repository) *Deleter {
	return &Deleter{repo: repo}
}

func (e *Deleter) Delete(ctx context.Context, fio types.FIO) error {
	return e.repo.Delete(ctx, fio)
}
