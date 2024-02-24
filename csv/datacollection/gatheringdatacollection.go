package datacollection

import (
	"github.com/level-5-pidgey/MarketMoogle/domain"
)

type GatheringDataCollection struct {
	GatheringItems      *map[int]domain.GatheringItem
	GatheringPointBases *map[int]domain.GatheringPointBase
	GatheringItemLevels *map[int]domain.GatheringItemLevel
	GatheringTypes      *map[int]domain.GatheringType
	GatheringPoints     *map[int][]domain.GatheringPoint
}
