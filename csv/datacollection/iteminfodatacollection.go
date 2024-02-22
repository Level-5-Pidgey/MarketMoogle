package datacollection

import csvType "github.com/level-5-pidgey/MarketMoogleApi/csv/types"

type ItemInfoDataCollection struct {
	Items                *map[int]csvType.Item
	Currencies           *map[int]csvType.Item
	ClassJobCategories   *map[int]csvType.ClassJobCategory
	ItemUiCategories     *map[int]csvType.ItemUiCategory
	ItemSearchCategories *map[int]csvType.ItemSearchCategory
	GilShopItems         *map[int]csvType.GilShopItem
	GcScripShopItem      *map[int]csvType.GcScripShopItem
}
