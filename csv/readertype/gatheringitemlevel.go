package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type GatheringItemLevel struct {
	Key   int
	Level int
	Stars int
}

func (g GatheringItemLevel) CreateFromCsvRow(record []string) (*GatheringItemLevel, error) {
	return &GatheringItemLevel{
		Key:   util.SafeStringToInt(record[0]),
		Level: util.SafeStringToInt(record[1]),
		Stars: util.SafeStringToInt(record[2]),
	}, nil
}

func (g GatheringItemLevel) GetKey() int {
	return g.Key
}
