package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type ClassJobCategory struct {
	Id          int
	JobCategory string
}

func (r ClassJobCategory) CreateFromCsvRow(record []string) (*ClassJobCategory, error) {
	return &ClassJobCategory{
		Id:          util.SafeStringToInt(record[0]),
		JobCategory: record[1],
	}, nil
}

func (r ClassJobCategory) GetKey() int {
	return r.Id
}
