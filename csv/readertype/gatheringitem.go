package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type GatheringItem struct {
	Key                   int
	ItemId                int
	GatheringItemLevelKey int
	IsHidden              bool
}

func (g GatheringItem) CreateFromCsvRow(record []string) (*GatheringItem, error) {
	itemId := util.SafeStringToInt(record[1])

	if itemId == 0 {
		return nil, nil
	}

	return &GatheringItem{
		Key:                   util.SafeStringToInt(record[0]),
		ItemId:                itemId,
		GatheringItemLevelKey: util.SafeStringToInt(record[2]),
		IsHidden:              util.SafeStringToBool(record[6]),
	}, nil
}

func (g GatheringItem) GetKey() int {
	return g.ItemId
}
