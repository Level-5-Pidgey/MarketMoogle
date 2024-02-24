package datacollection

import csvType "github.com/level-5-pidgey/MarketMoogle/csv/types"

type GatheringDataCollection struct {
	GatheringItems      *map[int]csvType.GatheringItem
	GatheringPointBases *map[int]csvType.GatheringPointBase
	GatheringItemLevels *map[int]csvType.GatheringItemLevel
	GatheringTypes      *map[int]csvType.GatheringType
	GatheringPoints     *map[int][]csvType.GatheringPoint
}
