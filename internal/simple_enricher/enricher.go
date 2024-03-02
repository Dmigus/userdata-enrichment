package usecases

import (
	"context"
	"errors"
	"sync"
	"userdata_enrichment/internal/usecases"
)

type AgeComputer interface {
	Get(context.Context, usecases.Key) (usecases.AgeType, error)
}

type SexComputer interface {
	Get(context.Context, usecases.Key) (usecases.SexType, error)
}

type NationalityComputer interface {
	Get(context.Context, usecases.Key) (usecases.NationalityType, error)
}

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

func (en *Enricher) Enrich(ctx context.Context, k usecases.Key) (usecases.Record, error) {
	wg := sync.WaitGroup{}
	wg.Add(3)
	var age usecases.AgeType
	var ageErr error
	go func() {
		defer wg.Done()
		age, ageErr = en.ageComp.Get(ctx, k)
	}()

	var sex usecases.SexType
	var sexErr error
	go func() {
		defer wg.Done()
		sex, sexErr = en.sexComp.Get(ctx, k)
	}()

	var nat usecases.NationalityType
	var natErr error
	go func() {
		defer wg.Done()
		nat, natErr = en.natComp.Get(ctx, k)
	}()

	wg.Wait()
	allErrs := errors.Join(ageErr, sexErr, natErr)
	if allErrs != nil {
		return usecases.Record{}, allErrs
	}
	return usecases.Record{
		Key:         k,
		Age:         age,
		Sex:         sex,
		Nationality: nat,
	}, nil
}
