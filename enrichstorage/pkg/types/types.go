package types

import (
	"fmt"
	"go.uber.org/zap"
)

type Age = int
type Sex = string
type Nationality = string

type FIO struct {
	name, surname, patronymic string
}

func NewFIO(name, surname, patronymic string) (FIO, error) {
	if len(name) == 0 || len(surname) == 0 {
		return FIO{}, fmt.Errorf("incorrect fio")
	}
	return FIO{name, surname, patronymic}, nil
}

func (fio FIO) hasPatronymic() bool {
	return len(fio.patronymic) > 0
}

func (fio FIO) Name() string {
	return fio.name
}

func (fio FIO) Surname() string {
	return fio.surname
}

func (fio FIO) Patronymic() string {
	return fio.patronymic
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
	Fio         FIO
	Age         Age
	Sex         Sex
	Nationality Nationality
}
