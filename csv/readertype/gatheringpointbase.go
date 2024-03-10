package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type GatheringPointBase struct {
	Key                 int
	GatheringTypeKey    int
	GatheringPointLevel int
	GatheringItemKeys   []int
}

func (g GatheringPointBase) CreateFromCsvRow(record []string) (*GatheringPointBase, error) {
	itemIds := make([]int, 0, 7)

	// Row index 3-10 contain items that can drop from this gathering point
	for i := 3; i < 10; i++ {
		itemId := util.SafeStringToInt(record[i])

		if itemId != 0 {
			itemIds = append(itemIds, itemId)
		}
	}

	if len(itemIds) == 0 {
		return nil, nil
	}

	return &GatheringPointBase{
		Key:                 util.SafeStringToInt(record[0]),
		GatheringTypeKey:    util.SafeStringToInt(record[1]),
		GatheringPointLevel: util.SafeStringToInt(record[2]),
		GatheringItemKeys:   itemIds,
	}, nil
}

func (g GatheringPointBase) GetKey() int {
	return g.Key
}
