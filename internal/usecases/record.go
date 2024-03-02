package usecases

type AgeType int
type SexType string
type NationalityType string

type Record struct {
	Key         Key
	Age         AgeType
	Sex         SexType
	Nationality NationalityType
}

type Key struct {
	name, surname, patronymic string
}
