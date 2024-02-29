package usecases

type Age int
type Sex string
type Nationality string

type Record struct {
	k           Key
	age         Age
	sex         Sex
	nationality Nationality
}

type Key struct {
	name, surname, patronymic string
}
