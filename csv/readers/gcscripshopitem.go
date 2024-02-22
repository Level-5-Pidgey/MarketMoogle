package csv

import (
	"github.com/level-5-pidgey/MarketMoogleApi/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogleApi/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogleApi/csv/types"
)

func NewGcScripShopItemReader() *csv.UngroupedXivApiCsvReader[csvType.GcScripShopItem] {
	csvColumns := map[string]int{
		"key":          0,
		"itemId":       1,
		"rankRequired": 2,
		"cost":         3,
	}

	return &csv.UngroupedXivApiCsvReader[csvType.GcScripShopItem]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.GcScripShopItem]{
			FileName:   "GCScripShopItem",
			RowsToSkip: 4,
			ProcessRow: func(record []string) (*csvType.GcScripShopItem, error) {
				return &csvType.GcScripShopItem{
					Key:                      csv.SafeStringToInt(record[csvColumns["key"]]),
					ItemId:                   csv.SafeStringToInt(record[csvColumns["itemId"]]),
					GrandCompanyRankRequired: csv.SafeStringToInt(record[csvColumns["rankRequired"]]),
					AmountRequired:           csv.SafeStringToInt(record[csvColumns["cost"]]),
				}, nil
			},
		},
	}
}
