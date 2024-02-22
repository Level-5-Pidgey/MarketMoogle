package csv

import (
	"github.com/level-5-pidgey/MarketMoogleApi/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogleApi/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogleApi/csv/types"
)

func NewClassJobCategoryReader() *csv.UngroupedXivApiCsvReader[csvType.ClassJobCategory] {
	csvColumns := map[string]int{
		"id":   0,
		"name": 1,
	}

	return &csv.UngroupedXivApiCsvReader[csvType.ClassJobCategory]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.ClassJobCategory]{
			FileName:   "ClassJobCategory",
			RowsToSkip: 4,
			ProcessRow: func(record []string) (*csvType.ClassJobCategory, error) {
				return &csvType.ClassJobCategory{
					Id:          csv.SafeStringToInt(record[csvColumns["id"]]),
					JobCategory: record[csvColumns["name"]],
				}, nil
			},
		},
	}
}
