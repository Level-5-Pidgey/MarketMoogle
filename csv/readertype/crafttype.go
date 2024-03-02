package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type CraftType struct {
	Key  int
	Name string
}

func (c CraftType) CreateFromCsvRow(record []string) (*CraftType, error) {
	return &CraftType{
		Key:  util.SafeStringToInt(record[0]),
		Name: record[3],
	}, nil
}

func (c CraftType) GetKey() int {
	return c.Key
}
