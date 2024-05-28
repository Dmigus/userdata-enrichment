package types

import "go.uber.org/zap"

type Age = int
type Sex = string
type Nationality = string

type FIO struct {
	name, surname, patronymic string
}

func (fio FIO) hasPatronymic() bool {
	return len(fio.patronymic) > 0
}

func FioToZaFields(fio FIO) []zap.Field {
	fs := make([]zap.Field, 0, 3)
	fs = append(fs, zap.String("name", fio.name))
	fs = append(fs, zap.String("surname", fio.surname))
	if fio.hasPatronymic() {
		fs = append(fs, zap.String("patronymic", fio.patronymic))
	}
	return fs
}

type EnrichedRecord struct {
	Key         FIO
	Age         Age
	Sex         Sex
	Nationality Nationality
}
