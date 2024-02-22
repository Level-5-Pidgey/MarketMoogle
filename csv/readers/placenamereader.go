package csv

import (
	"github.com/level-5-pidgey/MarketMoogleApi/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogleApi/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogleApi/csv/types"
)

func NewPlaceNameReader() *csv.UngroupedXivApiCsvReader[csvType.PlaceName] {
	csvColumns := map[string]int{
		"key":  0,
		"name": 1,
	}

	return &csv.UngroupedXivApiCsvReader[csvType.PlaceName]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.PlaceName]{
			FileName:   "PlaceName",
			RowsToSkip: 4,
			ProcessRow: func(record []string) (*csvType.PlaceName, error) {
				return &csvType.PlaceName{
					Key:  csv.SafeStringToInt(record[csvColumns["key"]]),
					Name: record[csvColumns["name"]],
				}, nil
			},
		},
	}
}
