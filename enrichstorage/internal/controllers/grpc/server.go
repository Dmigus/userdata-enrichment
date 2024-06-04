package grpc

import (
	"context"
	"enrichstorage/internal/controllers/grpc/converters"
	v1 "enrichstorage/internal/controllers/grpc/protoc"
	"enrichstorage/internal/service/enrichstorage/update"
	"enrichstorage/pkg/types"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type (
	PresenceChecker interface {
		IsFIOPresents(ctx context.Context, fio types.FIO) (bool, error)
	}
	Updater interface {
		Update(ctx context.Context, rec update.Request) error
	}
	Server struct {
		v1.UnimplementedEnrichStorageServer
		prChecker PresenceChecker
		updater   Updater
	}
)

func NewServer(prChecker PresenceChecker, updater Updater) *Server {
	return &Server{prChecker: prChecker, updater: updater}
}

func (s *Server) IsFIOPresents(ctx context.Context, fioDTO *v1.FIO) (*wrapperspb.BoolValue, error) {
	fio, err := converters.DtoToFIO(fioDTO)
	if err != nil {
		return nil, err
	}
	isPresents, err := s.prChecker.IsFIOPresents(ctx, fio)
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
	enrichReq := enrichedToUpdateRequest(enriched)
	err = s.updater.Update(ctx, enrichReq)
	return &emptypb.Empty{}, err
}

func enrichedToUpdateRequest(enriched types.EnrichedRecord) update.Request {
	return update.Request{
		Fio:                 enriched.Fio,
		NewAge:              enriched.Age,
		NewNat:              enriched.Nationality,
		NewSex:              enriched.Sex,
		SexPresents:         true,
		AgePresents:         true,
		NationalityPresents: true,
	}
}
