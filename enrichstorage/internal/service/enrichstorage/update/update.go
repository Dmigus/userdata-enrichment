package update

import (
	"context"
	"enrichstorage/pkg/types"
)

type (
	Request struct {
		Fio                                           types.FIO
		AgePresents, SexPresents, NationalityPresents bool
		NewAge                                        types.Age
		NewSex                                        types.Sex
		NewNat                                        types.Nationality
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
