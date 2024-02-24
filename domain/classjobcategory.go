package domain

import (
	"github.com/level-5-pidgey/MarketMoogle/util"
)

type ClassJobCategory struct {
	Id          int
	JobCategory string
}

func (r ClassJobCategory) GetKey() int {
	return r.Id
}

func (r ClassJobCategory) CreateFromCsvRow(record []string) (ReaderType, error) {
	return ClassJobCategory{
		Id:          util.SafeStringToInt(record[0]),
		JobCategory: record[1],
	}, nil
}
