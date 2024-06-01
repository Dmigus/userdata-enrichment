package converters

import (
	v1 "enricher/internal/providers/storage/protoc"
	"enrichstorage/pkg/types"
)

func FioToDTO(fio types.FIO) *v1.FIO {
	return &v1.FIO{
		Name:       fio.Name(),
		Surname:    fio.Surname(),
		Patronymic: fio.Patronymic(),
	}
}
