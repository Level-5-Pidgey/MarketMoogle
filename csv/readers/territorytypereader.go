package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/domain"
	"github.com/level-5-pidgey/MarketMoogle/util"
)

func NewTerritoryTypeReader() *csv.UngroupedXivApiCsvReader[csvType.TerritoryType] {
	csvColumns := map[string]int{
		"key":         0,
		"placeRegion": 4,
		"placeZone":   5,
		"placeName":   6,
		"map":         7,
	}

	return &csv.UngroupedXivApiCsvReader[csvType.TerritoryType]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.TerritoryType]{
			FileName:   "TerritoryType",
			RowsToSkip: 4,
			ProcessRow: func(record []string) (*csvType.TerritoryType, error) {
				return &csvType.TerritoryType{
					Key:      util.SafeStringToInt(record[csvColumns["key"]]),
					RegionId: util.SafeStringToInt(record[csvColumns["placeRegion"]]),
					PlaceId:  util.SafeStringToInt(record[csvColumns["placeName"]]),
					MapId:    util.SafeStringToInt(record[csvColumns["map"]]),
				}, nil
			},
		},
	}
}
