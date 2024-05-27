package delete

import (
	"bff/internal/service"
	"context"
)

type (
	repo interface {
		Delete(ctx context.Context, key service.FIO) error
	}
	Deleter struct {
		repo repo
	}
)

func NewDeleter(repo repo) *Deleter {
	return &Deleter{repo: repo}
}

func (e *Deleter) Delete(ctx context.Context, fio service.FIO) error {
	return e.repo.Delete(ctx, fio)
}
