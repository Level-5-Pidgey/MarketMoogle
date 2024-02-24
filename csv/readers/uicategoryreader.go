package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/domain"
	"github.com/level-5-pidgey/MarketMoogle/util"
)

func NewItemUiCategoryReader() *csv.UngroupedXivApiCsvReader[csvType.ItemUiCategory] {
	csvColumns := map[string]int{
		"id":     0,
		"name":   1,
		"iconId": 2,
	}

	return &csv.UngroupedXivApiCsvReader[csvType.ItemUiCategory]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.ItemUiCategory]{
			FileName:   "ItemUICategory",
			RowsToSkip: 4,
			ProcessRow: func(record []string) (*csvType.ItemUiCategory, error) {
				return &csvType.ItemUiCategory{
					Id:     util.SafeStringToInt(record[csvColumns["id"]]),
					Name:   record[csvColumns["name"]],
					IconId: util.SafeStringToInt(record[csvColumns["iconId"]]),
				}, nil
			},
		},
	}
}
