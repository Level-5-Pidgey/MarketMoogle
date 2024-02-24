package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/csv/types"
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
					Key:      csv.SafeStringToInt(record[csvColumns["key"]]),
					RegionId: csv.SafeStringToInt(record[csvColumns["placeRegion"]]),
					PlaceId:  csv.SafeStringToInt(record[csvColumns["placeName"]]),
					MapId:    csv.SafeStringToInt(record[csvColumns["map"]]),
				}, nil
			},
		},
	}
}
