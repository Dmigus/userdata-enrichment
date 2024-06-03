package grpc

import (
	"context"
	"enrichstorage/internal/controllers/grpc/converters"
	v1 "enrichstorage/internal/controllers/grpc/protoc"
	"enrichstorage/pkg/types"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type (
	service interface {
		IsFIOPresents(ctx context.Context, fio types.FIO) (bool, error)
		Update(ctx context.Context, rec types.EnrichedRecord) error
	}
	Server struct {
		v1.UnimplementedEnrichStorageServer
		service service
	}
)

func NewServer(service service) *Server {
	return &Server{service: service}
}

func (s *Server) IsFIOPresents(ctx context.Context, fioDTO *v1.FIO) (*wrapperspb.BoolValue, error) {
	fio, err := converters.DtoToFIO(fioDTO)
	if err != nil {
		return nil, err
	}
	isPresents, err := s.service.IsFIOPresents(ctx, fio)
	if err != nil {
		return nil, err
	}
	return &wrapperspb.BoolValue{Value: isPresents}, nil
}

func (s *Server) Update(ctx context.Context, enrichedDTO *v1.Enriched) (*emptypb.Empty, error) {
	enriched, err := converters.DtoToEnriched(enrichedDTO)
	if err != nil {
		return nil, err
	}
	err = s.service.Update(ctx, enriched)
	return &emptypb.Empty{}, err
}
