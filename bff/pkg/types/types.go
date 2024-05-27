package types

type Age = int
type Sex = string
type Nationality = string

type FIO struct {
	name, surname, patronymic string
}

type EnrichedRecord struct {
	Key         FIO
	Age         Age
	Sex         Sex
	Nationality Nationality
}
