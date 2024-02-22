package csv

import (
	"github.com/level-5-pidgey/MarketMoogleApi/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogleApi/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogleApi/csv/types"
)

func NewItemSearchCategoryReader() *csv.UngroupedXivApiCsvReader[csvType.ItemSearchCategory] {
	csvColumns := map[string]int{
		"id":       0,
		"name":     1,
		"icon":     2,
		"category": 3,
		"order":    4,
		"classJob": 5,
	}

	return &csv.UngroupedXivApiCsvReader[csvType.ItemSearchCategory]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.ItemSearchCategory]{
			FileName:   "ItemSearchCategory",
			RowsToSkip: 3,
			ProcessRow: func(record []string) (*csvType.ItemSearchCategory, error) {
				return &csvType.ItemSearchCategory{
					Key:           csv.SafeStringToInt(record[csvColumns["id"]]),
					Name:          record[csvColumns["name"]],
					IconId:        csv.SafeStringToInt(record[csvColumns["icon"]]),
					CategoryValue: csv.SafeStringToInt(record[csvColumns["category"]]),
					ClassJobId:    csv.SafeStringToInt(record[csvColumns["classJob"]]),
				}, nil
			},
		},
	}
}
