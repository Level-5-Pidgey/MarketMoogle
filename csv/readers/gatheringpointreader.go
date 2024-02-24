package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/csv/types"
)

func NewGatheringPointReader() *csv.GroupedXivApiCsvReader[csvType.GatheringPoint] {
	csvColumns := map[string]int{
		"key":                0,
		"gatheringType":      1,
		"gatheringPointBase": 3,
		"territoryType":      7,
		"placeName":          8,
	}

	return &csv.GroupedXivApiCsvReader[csvType.GatheringPoint]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.GatheringPoint]{
			FileName:   "GatheringPoint",
			RowsToSkip: 4,
			ProcessRow: func(record []string) (*csvType.GatheringPoint, error) {
				return &csvType.GatheringPoint{
					Key:                  csv.SafeStringToInt(record[csvColumns["key"]]),
					GatheringTypeId:      csv.SafeStringToInt(record[csvColumns["gatheringType"]]),
					GatheringPointBaseId: csv.SafeStringToInt(record[csvColumns["gatheringPointBase"]]),
					TerritoryTypeId:      csv.SafeStringToInt(record[csvColumns["territoryType"]]),
					PlaceNameId:          csv.SafeStringToInt(record[csvColumns["placeName"]]),
				}, nil
			},
		},
	}
}
