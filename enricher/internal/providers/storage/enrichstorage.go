package storage

import (
	"context"
	v1 "enricher/internal/providers/storage/protoc"
	"enrichstorage/pkg/types"
)

type EnrichStorage struct {
	client v1.EnrichStorageClient
}

func (e EnrichStorage) IsFIOPresents(ctx context.Context, fio types.FIO) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (e EnrichStorage) Update(ctx context.Context, rec types.EnrichedRecord) error {
	//TODO implement me
	panic("implement me")
}
