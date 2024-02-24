package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/domain"
	"github.com/level-5-pidgey/MarketMoogle/util"
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
					Key:           util.SafeStringToInt(record[csvColumns["id"]]),
					Name:          record[csvColumns["name"]],
					IconId:        util.SafeStringToInt(record[csvColumns["icon"]]),
					CategoryValue: util.SafeStringToInt(record[csvColumns["category"]]),
					ClassJobId:    util.SafeStringToInt(record[csvColumns["classJob"]]),
				}, nil
			},
		},
	}
}
