package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/domain"
	"github.com/level-5-pidgey/MarketMoogle/util"
)

func NewGatheringPointBaseReader() *csv.UngroupedXivApiCsvReader[csvType.GatheringPointBase] {
	csvColumns := map[string]int{
		"key":            0,
		"gatheringType":  1,
		"gatheringLevel": 2,
	}

	itemColumns := map[string]int{
		"item1": 3,
		"item2": 4,
		"item3": 5,
		"item4": 6,
		"item5": 7,
		"item6": 8,
		"item7": 9,
		"item8": 10,
	}

	return &csv.UngroupedXivApiCsvReader[csvType.GatheringPointBase]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.GatheringPointBase]{
			FileName:   "GatheringPointBase",
			RowsToSkip: 4,
			ProcessRow: func(record []string) (*csvType.GatheringPointBase, error) {
				itemIds := make([]int, 0)

				for _, column := range itemColumns {
					itemId := util.SafeStringToInt(record[column])

					if itemId != 0 {
						itemIds = append(itemIds, itemId)
					}
				}

				if len(itemIds) == 0 {
					return nil, nil
				}

				return &csvType.GatheringPointBase{
					Key:                 util.SafeStringToInt(record[csvColumns["key"]]),
					GatheringTypeKey:    util.SafeStringToInt(record[csvColumns["gatheringType"]]),
					GatheringPointLevel: util.SafeStringToInt(record[csvColumns["gatheringLevel"]]),
					GatheringItemKeys:   itemIds,
				}, nil
			},
		},
	}
}
