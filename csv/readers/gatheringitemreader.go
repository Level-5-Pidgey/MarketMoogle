package csv

import (
	"github.com/level-5-pidgey/MarketMoogleApi/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogleApi/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogleApi/csv/types"
)

func NewItemGatheringItemReader() *csv.UngroupedXivApiCsvReader[csvType.GatheringItem] {
	csvColumns := map[string]int{
		"key":                0,
		"itemId":             1,
		"gatheringItemLevel": 2,
		"hidden":             6,
	}

	return &csv.UngroupedXivApiCsvReader[csvType.GatheringItem]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.GatheringItem]{
			FileName:   "GatheringItem",
			RowsToSkip: 4,
			ProcessRow: func(record []string) (*csvType.GatheringItem, error) {
				itemId := csv.SafeStringToInt(record[csvColumns["itemId"]])

				if itemId == 0 {
					return nil, nil
				}

				return &csvType.GatheringItem{
					Key:                   csv.SafeStringToInt(record[csvColumns["key"]]),
					ItemId:                itemId,
					GatheringItemLevelKey: csv.SafeStringToInt(record[csvColumns["gatheringItemLevel"]]),
					IsHidden:              csv.SafeStringToBool(record[csvColumns["hidden"]]),
				}, nil
			},
		},
	}
}
