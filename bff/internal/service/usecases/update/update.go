package update

import (
	"bff/internal/service"
	"context"
)

type (
	Request struct {
		fio           service.FIO
		age, sex, nat bool
		newAge        service.Age
		newSex        service.Sex
		newNat        service.Nationality
	}
	repo interface {
		Update(ctx context.Context, req Request) error
	}
	Updater struct {
		repo repo
	}
)

func NewUpdater(repo repo) *Updater {
	return &Updater{repo: repo}
}

func (e *Updater) Update(ctx context.Context, req Request) error {
	return e.repo.Update(ctx, req)
}
