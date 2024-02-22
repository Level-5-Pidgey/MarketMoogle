package csv

import (
	"github.com/level-5-pidgey/MarketMoogleApi/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogleApi/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogleApi/csv/types"
)

func NewCraftTypeReader() *csv.UngroupedXivApiCsvReader[csvType.CraftType] {
	csvColumns := map[string]int{
		"key":  0,
		"name": 3,
	}

	return &csv.UngroupedXivApiCsvReader[csvType.CraftType]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.CraftType]{
			FileName:   "CraftType",
			RowsToSkip: 3,
			ProcessRow: func(record []string) (*csvType.CraftType, error) {
				return &csvType.CraftType{
					Key:  csv.SafeStringToInt(record[csvColumns["key"]]),
					Name: record[csvColumns["name"]],
				}, nil
			},
		},
	}
}
