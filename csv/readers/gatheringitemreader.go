package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/domain"
	"github.com/level-5-pidgey/MarketMoogle/util"
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
				itemId := util.SafeStringToInt(record[csvColumns["itemId"]])

				if itemId == 0 {
					return nil, nil
				}

				return &csvType.GatheringItem{
					Key:                   util.SafeStringToInt(record[csvColumns["key"]]),
					ItemId:                itemId,
					GatheringItemLevelKey: util.SafeStringToInt(record[csvColumns["gatheringItemLevel"]]),
					IsHidden:              util.SafeStringToBool(record[csvColumns["hidden"]]),
				}, nil
			},
		},
	}
}
