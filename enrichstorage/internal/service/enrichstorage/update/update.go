package update

import (
	"context"
	"enrichstorage/pkg/types"
)

type (
	Request struct {
		fio           types.FIO
		age, sex, nat bool
		newAge        types.Age
		newSex        types.Sex
		newNat        types.Nationality
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
