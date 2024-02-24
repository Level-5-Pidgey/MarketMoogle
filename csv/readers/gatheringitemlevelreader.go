package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/csv/types"
)

func NewGatheringItemLevelReader() *csv.UngroupedXivApiCsvReader[csvType.GatheringItemLevel] {
	csvColumns := map[string]int{
		"key":   0,
		"level": 1,
		"stars": 2,
	}

	return &csv.UngroupedXivApiCsvReader[csvType.GatheringItemLevel]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.GatheringItemLevel]{
			FileName:   "GatheringItemLevelConvertTable",
			RowsToSkip: 4,
			ProcessRow: func(record []string) (*csvType.GatheringItemLevel, error) {
				return &csvType.GatheringItemLevel{
					Key:   csv.SafeStringToInt(record[csvColumns["key"]]),
					Level: csv.SafeStringToInt(record[csvColumns["level"]]),
					Stars: csv.SafeStringToInt(record[csvColumns["stars"]]),
				}, nil
			},
		},
	}
}
