package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type GilShopItem struct {
	Key    int
	ItemId int
}

func (g GilShopItem) CreateFromCsvRow(record []string) (*GilShopItem, error) {
	return &GilShopItem{
		Key:    util.SafeStringToInt(record[0]),
		ItemId: util.SafeStringToInt(record[1]),
	}, nil
}

func (g GilShopItem) GetKey() int {
	return g.ItemId
}
