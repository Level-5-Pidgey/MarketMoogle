package datacollection

import csvType "github.com/level-5-pidgey/MarketMoogle/csv/types"

type PlaceDataCollection struct {
	PlaceNames     *map[int]csvType.PlaceName
	TerritoryTypes *map[int]csvType.TerritoryType
}
