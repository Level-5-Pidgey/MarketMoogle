package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type ItemSearchCategory struct {
	Key           int
	Name          string
	IconId        int
	CategoryValue int
	ClassJobId    int
}

func (i ItemSearchCategory) CreateFromCsvRow(record []string) (*ItemSearchCategory, error) {
	return &ItemSearchCategory{
		Key:           util.SafeStringToInt(record[0]),
		Name:          record[1],
		IconId:        util.SafeStringToInt(record[2]),
		CategoryValue: util.SafeStringToInt(record[3]),
		ClassJobId:    util.SafeStringToInt(record[5]),
	}, nil
}

func (i ItemSearchCategory) GetKey() int {
	return i.Key
}
