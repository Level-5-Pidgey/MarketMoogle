package csv

import (
	"github.com/level-5-pidgey/MarketMoogleApi/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogleApi/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogleApi/csv/types"
)

func NewGilShopItemReader() *csv.UngroupedXivApiCsvReader[csvType.GilShopItem] {
	csvColumns := map[string]int{
		"key":    0,
		"itemId": 1,
	}

	return &csv.UngroupedXivApiCsvReader[csvType.GilShopItem]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.GilShopItem]{
			FileName:   "GilShopItem",
			RowsToSkip: 3,
			ProcessRow: func(record []string) (*csvType.GilShopItem, error) {
				return &csvType.GilShopItem{
					Key:    csv.SafeStringToInt(record[csvColumns["key"]]),
					ItemId: csv.SafeStringToInt(record[csvColumns["itemId"]]),
				}, nil
			},
		},
	}
}
