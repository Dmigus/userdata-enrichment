package usecases

import (
	"context"
	"errors"
	"sync"
)

type AgeComputer interface {
	Get(context.Context, service.Key) (service.AgeType, error)
}

type SexComputer interface {
	Get(context.Context, service.Key) (service.SexType, error)
}

type NationalityComputer interface {
	Get(context.Context, service.Key) (service.NationalityType, error)
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

func (en *Enricher) Enrich(ctx context.Context, k service.Key) (service.Record, error) {
	wg := sync.WaitGroup{}
	wg.Add(3)
	var age service.AgeType
	var ageErr error
	go func() {
		defer wg.Done()
		age, ageErr = en.ageComp.Get(ctx, k)
	}()

	var sex service.SexType
	var sexErr error
	go func() {
		defer wg.Done()
		sex, sexErr = en.sexComp.Get(ctx, k)
	}()

	var nat service.NationalityType
	var natErr error
	go func() {
		defer wg.Done()
		nat, natErr = en.natComp.Get(ctx, k)
	}()

	wg.Wait()
	allErrs := errors.Join(ageErr, sexErr, natErr)
	if allErrs != nil {
		return service.Record{}, allErrs
	}
	return service.Record{
		Key:         k,
		Age:         age,
		Sex:         sex,
		Nationality: nat,
	}, nil
}
