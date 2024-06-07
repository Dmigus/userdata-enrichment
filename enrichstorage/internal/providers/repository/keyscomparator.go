package repository

import (
	"cmp"
	"enrichstorage/pkg/types"
)

type FioComparator struct{}

func (f *FioComparator) Cmp(a, b types.FIO) int {
	if res := cmp.Compare(a.Surname(), b.Surname()); res != 0 {
		return res
	}
	if res := cmp.Compare(a.Name(), b.Name()); res != 0 {
		return res
	}
	return cmp.Compare(a.Patronymic(), b.Patronymic())
}
