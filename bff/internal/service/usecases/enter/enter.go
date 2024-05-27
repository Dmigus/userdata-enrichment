package enter

import (
	"bff/internal/service"
	"context"
)

type (
	repo interface {
		Create(ctx context.Context, key service.FIO) error
	}
	handleQueue interface {
		Push(ctx context.Context, key service.FIO) error
	}
	Enterer struct {
		repo  repo
		queue handleQueue
	}
)

func NewEnterer(repo repo, queue handleQueue) *Enterer {
	return &Enterer{repo: repo, queue: queue}
}

func (e *Enterer) Enter(ctx context.Context, fio service.FIO) error {
	if err := e.queue.Push(ctx, fio); err != nil {
		return err
	}
	return e.repo.Create(ctx, fio)
}
