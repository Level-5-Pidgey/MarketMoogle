package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/domain"
	"github.com/level-5-pidgey/MarketMoogle/util"
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
					Key:    util.SafeStringToInt(record[csvColumns["key"]]),
					ItemId: util.SafeStringToInt(record[csvColumns["itemId"]]),
				}, nil
			},
		},
	}
}
