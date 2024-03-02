package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type GatheringType struct {
	Key    int
	Name   string
	IconId int
}

func (g GatheringType) CreateFromCsvRow(record []string) (*GatheringType, error) {
	return &GatheringType{
		Key:    util.SafeStringToInt(record[0]),
		Name:   record[1],
		IconId: util.SafeStringToInt(record[2]),
	}, nil
}

func (g GatheringType) GetKey() int {
	return g.Key
}
