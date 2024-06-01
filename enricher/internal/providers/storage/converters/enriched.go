package converters

import (
	v1 "enricher/internal/providers/storage/protoc"
	"enrichstorage/pkg/types"
)

func EnrichedToDto(en types.EnrichedRecord) *v1.Enriched {
	fioDTO := FioToDTO(en.Key)
	return &v1.Enriched{
		Fio:         fioDTO,
		Age:         int32(en.Age),
		Sex:         en.Sex,
		Nationality: en.Nationality,
	}
}
