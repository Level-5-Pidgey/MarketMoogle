package datacollection

import csvType "github.com/level-5-pidgey/MarketMoogleApi/csv/types"

type PlaceDataCollection struct {
	PlaceNames     *map[int]csvType.PlaceName
	TerritoryTypes *map[int]csvType.TerritoryType
}
