package csv

type ClassJobCategory struct {
	Id          int
	JobCategory string
}

func (r ClassJobCategory) GetKey() int {
	return r.Id
}
