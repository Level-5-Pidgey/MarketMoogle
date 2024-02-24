package datacollection

import (
	"github.com/level-5-pidgey/MarketMoogle/domain"
)

type ItemInfoDataCollection struct {
	Items                *map[int]domain.Item
	Currencies           *map[int]domain.Item
	ClassJobCategories   *map[int]domain.ClassJobCategory
	ItemUiCategories     *map[int]domain.ItemUiCategory
	ItemSearchCategories *map[int]domain.ItemSearchCategory
	GilShopItems         *map[int]domain.GilShopItem
	GcScripShopItem      *map[int]domain.GcScripShopItem
}
