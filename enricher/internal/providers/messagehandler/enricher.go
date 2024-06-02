package messagehandler

import (
	"context"
	"enrichstorage/pkg/types"
	"golang.org/x/sync/errgroup"
)

type AgeComputer interface {
	Get(context.Context, types.FIO) (types.Age, error)
}

type SexComputer interface {
	Get(context.Context, types.FIO) (types.Sex, error)
}

type NationalityComputer interface {
	Get(context.Context, types.FIO) (types.Nationality, error)
}

// Enricher обогащает информацией по ФИО
type Enricher struct {
	ageComp AgeComputer
	sexComp SexComputer
	natComp NationalityComputer
}

func New(ageComp AgeComputer, sexComp SexComputer, natComp NationalityComputer) *Enricher {
	return &Enricher{
		ageComp,
		sexComp,
		natComp,
	}
}

func (en *Enricher) Enrich(ctx context.Context, k types.FIO) (types.EnrichedRecord, error) {
	eg, ctx := errgroup.WithContext(ctx)
	var age types.Age
	eg.Go(func() error {
		var err error
		age, err = en.ageComp.Get(ctx, k)
		return err
	})

	var sex types.Sex
	eg.Go(func() error {
		var err error
		sex, err = en.sexComp.Get(ctx, k)
		return err
	})

	var nat types.Nationality
	eg.Go(func() error {
		var err error
		nat, err = en.natComp.Get(ctx, k)
		return err
	})

	if err := eg.Wait(); err != nil {
		return types.EnrichedRecord{}, err
	}

	return types.EnrichedRecord{
		Fio:         k,
		Age:         age,
		Sex:         sex,
		Nationality: nat,
	}, nil
}
