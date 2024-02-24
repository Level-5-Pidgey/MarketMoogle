package domain

type ItemUiCategory struct {
	Id     int
	Name   string
	IconId int
}

func (r ItemUiCategory) GetKey() int {
	return r.Id
}
