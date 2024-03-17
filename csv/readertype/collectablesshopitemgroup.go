package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type CollectablesShopItemGroup struct {
	Key  int
	Name string
}

func (c CollectablesShopItemGroup) GetKey() int {
	return c.Key
}

func (c CollectablesShopItemGroup) CreateFromCsvRow(record []string) (*CollectablesShopItemGroup, error) {
	return &CollectablesShopItemGroup{
		Key:  util.SafeStringToInt(record[0]),
		Name: record[1],
	}, nil
}
