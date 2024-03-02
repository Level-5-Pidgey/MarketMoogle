package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type GcScripShopItem struct {
	Key                      int
	ItemId                   int
	GrandCompanyRankRequired int
	AmountRequired           int
}

func (g GcScripShopItem) CreateFromCsvRow(record []string) (*GcScripShopItem, error) {
	return &GcScripShopItem{
		Key:                      util.SafeStringToInt(record[0]),
		ItemId:                   util.SafeStringToInt(record[1]),
		GrandCompanyRankRequired: util.SafeStringToInt(record[2]),
		AmountRequired:           util.SafeStringToInt(record[3]),
	}, nil
}

func (g GcScripShopItem) GetKey() int {
	return g.ItemId
}
