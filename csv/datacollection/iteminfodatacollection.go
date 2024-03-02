package datacollection

import (
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
)

type ItemInfoDataCollection struct {
	Items                *map[int]*readertype.Item
	Currencies           *map[int]*readertype.Item
	ClassJobCategories   *map[int]*readertype.ClassJobCategory
	ItemUiCategories     *map[int]*readertype.ItemUiCategory
	ItemSearchCategories *map[int]*readertype.ItemSearchCategory
	GilShopItems         *map[int]*readertype.GilShopItem
	GcScripShopItem      *map[int]*readertype.GcScripShopItem
}
