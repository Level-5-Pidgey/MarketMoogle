package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/domain"
	"github.com/level-5-pidgey/MarketMoogle/util"
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
					Key:                      util.SafeStringToInt(record[csvColumns["key"]]),
					ItemId:                   util.SafeStringToInt(record[csvColumns["itemId"]]),
					GrandCompanyRankRequired: util.SafeStringToInt(record[csvColumns["rankRequired"]]),
					AmountRequired:           util.SafeStringToInt(record[csvColumns["cost"]]),
				}, nil
			},
		},
	}
}
