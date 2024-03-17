package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type CollectablesShopItem struct {
	Key         float64
	ItemId      int
	ItemGroup   int
	LevelMin    int
	LevelMax    int
	Stars       int
	RewardScrip int
}

func (c CollectablesShopItem) GetKey() int {
	return c.ItemId
}

func (c CollectablesShopItem) CreateFromCsvRow(record []string) (*CollectablesShopItem, error) {
	return &CollectablesShopItem{
		Key:         util.SafeStringToFloat(record[0]),
		ItemId:      util.SafeStringToInt(record[1]),
		ItemGroup:   util.SafeStringToInt(record[2]),
		LevelMin:    util.SafeStringToInt(record[3]),
		LevelMax:    util.SafeStringToInt(record[5]),
		Stars:       util.SafeStringToInt(record[6]),
		RewardScrip: util.SafeStringToInt(record[9]),
	}, nil
}
