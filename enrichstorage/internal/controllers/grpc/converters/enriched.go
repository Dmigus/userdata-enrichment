package converters

import (
	v1 "enrichstorage/internal/controllers/grpc/protoc"
	"enrichstorage/pkg/types"
)

func DtoToEnriched(dto *v1.Enriched) (types.EnrichedRecord, error) {
	fio, err := DtoToFIO(dto.GetFio())
	if err != nil {
		return types.EnrichedRecord{}, err
	}
	return types.EnrichedRecord{
		Fio:         fio,
		Age:         types.Age(dto.GetAge()),
		Sex:         dto.GetSex(),
		Nationality: dto.GetNationality(),
	}, nil
}
