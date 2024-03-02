package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type PlaceName struct {
	Key  int
	Name string
}

func (p PlaceName) CreateFromCsvRow(record []string) (*PlaceName, error) {
	return &PlaceName{
		Key:  util.SafeStringToInt(record[0]),
		Name: record[1],
	}, nil
}

func (p PlaceName) GetKey() int {
	return p.Key
}
