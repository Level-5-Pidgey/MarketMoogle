package datacollection

import "github.com/level-5-pidgey/MarketMoogle/csv/readertype"

type GatheringDataCollection struct {
	GatheringItems             *map[int]*readertype.GatheringItem
	GatheringPointBases        *map[int]*readertype.GatheringPointBase
	GatheringItemLevels        *map[int]*readertype.GatheringItemLevel
	GatheringTypes             *map[int]*readertype.GatheringType
	GatheringPoints            *map[int][]*readertype.GatheringPoint
	CollectablesShopItem       *map[int]*readertype.CollectablesShopItem
	CollectableShopRewardScrip *map[int]*readertype.CollectableShopRewardScrip
	CollectablesShopItemGroup  *map[int]*readertype.CollectablesShopItemGroup
}
