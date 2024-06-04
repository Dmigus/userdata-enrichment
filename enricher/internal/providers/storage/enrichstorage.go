package storage

import (
	"context"
	"enricher/internal/providers/storage/converters"
	v1 "enricher/internal/providers/storage/protoc"
	"enrichstorage/pkg/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EnrichStorage struct {
	client v1.EnrichStorageClient
}

func NewEnrichStorage(addr string) (*EnrichStorage, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	client := v1.NewEnrichStorageClient(conn)
	return &EnrichStorage{client: client}, nil
}

func (e *EnrichStorage) IsFIOPresents(ctx context.Context, fio types.FIO) (bool, error) {
	dto := converters.FioToDTO(fio)
	isPresents, err := e.client.IsFIOPresents(ctx, dto)
	if err != nil {
		return false, err
	}
	if isPresents == nil {
		return false, nil
	}
	return isPresents.Value, nil
}

func (e *EnrichStorage) Update(ctx context.Context, rec types.EnrichedRecord) error {
	dto := converters.EnrichedToDto(rec)
	_, err := e.client.Update(ctx, dto)
	return err
}
