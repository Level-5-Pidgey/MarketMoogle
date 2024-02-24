package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/domain"
	"github.com/level-5-pidgey/MarketMoogle/util"
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
					Key:                  util.SafeStringToInt(record[csvColumns["key"]]),
					GatheringTypeId:      util.SafeStringToInt(record[csvColumns["gatheringType"]]),
					GatheringPointBaseId: util.SafeStringToInt(record[csvColumns["gatheringPointBase"]]),
					TerritoryTypeId:      util.SafeStringToInt(record[csvColumns["territoryType"]]),
					PlaceNameId:          util.SafeStringToInt(record[csvColumns["placeName"]]),
				}, nil
			},
		},
	}
}
