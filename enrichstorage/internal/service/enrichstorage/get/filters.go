package get

import "enrichstorage/pkg/types"

type (
	NameFilter struct {
		Val string
	}
	SurnameFilter struct {
		Val string
	}
	PatronymicFilter struct {
		Val string
	}
	SexFilter struct {
		Val string
	}
	AgeFilter struct {
		LTE, GTE types.Age
	}
	NationalityFilter struct {
		Val string
	}
	Filters struct {
		nameFilter       *NameFilter
		surnameFilter    *SurnameFilter
		patronymicFilter *PatronymicFilter
		sexFilter        *SexFilter
		ageFilter        *AgeFilter
		natFilter        *NationalityFilter
	}
)

func (f *Filters) NameFilter() (*NameFilter, bool) {
	return f.nameFilter, f.nameFilter != nil
}

func (f *Filters) SetNameFilter(nameFilter *NameFilter) {
	f.nameFilter = nameFilter
}

func (f *Filters) SurnameFilter() (*SurnameFilter, bool) {
	return f.surnameFilter, f.surnameFilter != nil
}

func (f *Filters) SetSurnameFilter(surnameFilter *SurnameFilter) {
	f.surnameFilter = surnameFilter
}

func (f *Filters) PatronymicFilter() (*PatronymicFilter, bool) {
	return f.patronymicFilter, f.patronymicFilter != nil
}

func (f *Filters) SetPatronymicFilter(patronymicFilter *PatronymicFilter) {
	f.patronymicFilter = patronymicFilter
}

func (f *Filters) SexFilter() (*SexFilter, bool) {
	return f.sexFilter, f.sexFilter != nil
}

func (f *Filters) SetSexFilter(sexFilter *SexFilter) {
	f.sexFilter = sexFilter
}

func (f *Filters) AgeFilter() (*AgeFilter, bool) {
	return f.ageFilter, f.ageFilter != nil
}

func (f *Filters) SetAgeFilter(ageFilter *AgeFilter) {
	f.ageFilter = ageFilter
}

func (f *Filters) NatFilter() (*NationalityFilter, bool) {
	return f.natFilter, f.natFilter != nil
}

func (f *Filters) SetNatFilter(natFilter *NationalityFilter) {
	f.natFilter = natFilter
}

func NewFilters() *Filters {
	return &Filters{}
}
