package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/domain"
	"github.com/level-5-pidgey/MarketMoogle/util"
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
					Id:          util.SafeStringToInt(record[csvColumns["id"]]),
					JobCategory: record[csvColumns["name"]],
				}, nil
			},
		},
	}
}
