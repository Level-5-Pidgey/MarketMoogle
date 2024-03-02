package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type ItemUiCategory struct {
	Id     int
	Name   string
	IconId int
}

func (r ItemUiCategory) CreateFromCsvRow(record []string) (*ItemUiCategory, error) {
	return &ItemUiCategory{
		Id:     util.SafeStringToInt(record[0]),
		Name:   record[1],
		IconId: util.SafeStringToInt(record[2]),
	}, nil
}

func (r ItemUiCategory) GetKey() int {
	return r.Id
}
