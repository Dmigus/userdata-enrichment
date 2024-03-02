package enrich

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"userdata_enrichment/internal/usecases"
)

type Enricher interface {
	Enrich(ctx context.Context, k usecases.Key) (usecases.Record, error)
}

type PullQueue interface {
	Pull(ctx context.Context) (usecases.Key, error)
}

type KVStorage interface {
	EnrichIfPresent(context.Context, usecases.Record) error
	IsPresent(context.Context, usecases.Key) bool
}

type Transaction interface {
	KVStorage
	PullQueue
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type TransactionalStorage interface {
	BeginTx(ctx context.Context) (Transaction, error)
}

type Service struct {
	started atomic.Bool
	cf      context.CancelFunc
	enr     Enricher
	db      TransactionalStorage
}

func (es *Service) Run() {
	if !es.started.CompareAndSwap(false, true) {
		return
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	es.cf = cancelFunc
	go func() {
		for {
			err := es.oneCycle(ctx)
			if err != nil && !errors.Is(err, context.Canceled) {
				log.Println(err)
			}
		}
	}()
}

func (es *Service) oneCycle(ctx context.Context) error {
	tr, err := es.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tr.Rollback(ctx)
	keyToEnrich, err := tr.Pull(ctx)
	if err != nil {
		return err
	}
	if !tr.IsPresent(ctx, keyToEnrich) {
		return tr.Commit(ctx)
	}
	enriched, err := es.enr.Enrich(ctx, keyToEnrich)
	if err != nil {
		return err
	}
	if err = tr.EnrichIfPresent(ctx, enriched); err != nil {
		return err
	}
	return tr.Commit(ctx)
}

func (es *Service) Stop() {
	es.cf()
	es.started.Store(false)
}
