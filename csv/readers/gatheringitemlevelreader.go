package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/domain"
	"github.com/level-5-pidgey/MarketMoogle/util"
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
					Key:   util.SafeStringToInt(record[csvColumns["key"]]),
					Level: util.SafeStringToInt(record[csvColumns["level"]]),
					Stars: util.SafeStringToInt(record[csvColumns["stars"]]),
				}, nil
			},
		},
	}
}
