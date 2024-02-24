package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/domain"
	"github.com/level-5-pidgey/MarketMoogle/util"
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
					Key:  util.SafeStringToInt(record[csvColumns["key"]]),
					Name: record[csvColumns["name"]],
				}, nil
			},
		},
	}
}
