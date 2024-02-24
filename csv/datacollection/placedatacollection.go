package datacollection

import (
	"github.com/level-5-pidgey/MarketMoogle/domain"
)

type PlaceDataCollection struct {
	PlaceNames     *map[int]domain.PlaceName
	TerritoryTypes *map[int]domain.TerritoryType
}
