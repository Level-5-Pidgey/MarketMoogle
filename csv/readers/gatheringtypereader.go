package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/csv/types"
)

func NewItemGatheringTypeReader() *csv.UngroupedXivApiCsvReader[csvType.GatheringType] {
	csvColumns := map[string]int{
		"key":    0,
		"name":   1,
		"iconId": 2,
	}

	return &csv.UngroupedXivApiCsvReader[csvType.GatheringType]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.GatheringType]{
			FileName:   "GatheringType",
			RowsToSkip: 3,
			ProcessRow: func(record []string) (*csvType.GatheringType, error) {
				return &csvType.GatheringType{
					Key:    csv.SafeStringToInt(record[csvColumns["key"]]),
					Name:   record[csvColumns["name"]],
					IconId: csv.SafeStringToInt(record[csvColumns["iconId"]]),
				}, nil
			},
		},
	}
}
