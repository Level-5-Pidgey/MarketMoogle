package datacollection

import "github.com/level-5-pidgey/MarketMoogle/csv/readertype"

type PlaceDataCollection struct {
	PlaceNames     *map[int]*readertype.PlaceName
	TerritoryTypes *map[int]*readertype.TerritoryType
}
