package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/csv/types"
	"github.com/level-5-pidgey/MarketMoogle/util"
)

func NewWorldReader() *csv.UngroupedXivApiCsvReader[csvType.World] {
	csvColumns := map[string]int{
		"id":         0,
		"name":       2,
		"region":     3,
		"dataCenter": 5,
		"isPublic":   6,
	}

	return &csv.UngroupedXivApiCsvReader[csvType.World]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.World]{
			FileName:   "World",
			RowsToSkip: 11,
			ProcessRow: func(record []string) (*csvType.World, error) {
				isPublic := util.SafeStringToBool(record[csvColumns["isPublic"]])

				if !isPublic {
					return nil, nil
				}

				return &csvType.World{
					Key:        util.SafeStringToInt(record[csvColumns["id"]]),
					Name:       record[csvColumns["name"]],
					Region:     util.SafeStringToInt(record[csvColumns["region"]]),
					DataCenter: util.SafeStringToInt(record[csvColumns["dataCenter"]]),
					IsPublic:   isPublic,
				}, nil
			},
		},
	}
}
