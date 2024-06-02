package converters

import (
	v1 "enrichstorage/internal/controllers/grpc/protoc"
	"enrichstorage/pkg/types"
)

func DtoToFIO(dto *v1.FIO) (types.FIO, error) {
	return types.NewFIO(dto.GetName(), dto.GetSurname(), dto.GetPatronymic())
}
