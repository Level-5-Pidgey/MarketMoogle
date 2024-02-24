package csv

import (
	"errors"
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/csv/types"
	"strings"
)

func NewItemCsvReader() *csv.UngroupedXivApiCsvReader[csvType.Item] {
	csvColumns := map[string]int{
		"id":                 0,
		"description":        9,
		"name":               10,
		"iconId":             11,
		"itemLevel":          12,
		"rarity":             13,
		"uiCategory":         16,
		"searchCategory":     17,
		"sortCategory":       19,
		"stackSize":          21,
		"isUntradable":       23,
		"dungeonDrop":        25,
		"buyFromVendorPrice": 26,
		"sellToVendorPrice":  27,
		"canBeHq":            28,
		"canDesynth":         37,
		"alwaysCollectable":  39,
		"classJobCategory":   44,
		"isGlamourous":       91,
	}

	return &csv.UngroupedXivApiCsvReader[csvType.Item]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.Item]{
			FileName:   "Item",
			RowsToSkip: 2,
			ProcessRow: func(record []string) (*csvType.Item, error) {
				itemId := csv.SafeStringToInt(record[csvColumns["id"]])
				itemName := record[csvColumns["name"]]

				if isOldItem(itemId, itemName) {
					return nil, errors.New("dated item, skipped")
				}

				if itemName == "" {
					return nil, errors.New("item name is empty")
				}

				canBeTraded := !csv.SafeStringToBool(record[csvColumns["isUntradable"]])
				adjustedBuyPrice := 0
				if canBeTraded {
					adjustedBuyPrice = csv.SafeStringToInt(record[csvColumns["buyFromVendorPrice"]])
				}

				canDesynth := false
				desynthsTo := csv.SafeStringToInt(record[csvColumns["canDesynth"]])
				if desynthsTo > 0 {
					canDesynth = true
				}

				result := csvType.Item{
					Id:                 itemId,
					Name:               record[csvColumns["name"]],
					Description:        record[csvColumns["description"]],
					IconId:             csv.SafeStringToInt(record[csvColumns["iconId"]]),
					ItemLevel:          csv.SafeStringToInt(record[csvColumns["itemLevel"]]),
					Rarity:             csv.SafeStringToInt(record[csvColumns["rarity"]]),
					UiCategory:         csv.SafeStringToInt(record[csvColumns["uiCategory"]]),
					SearchCategory:     csv.SafeStringToInt(record[csvColumns["searchCategory"]]),
					SortCategory:       csv.SafeStringToInt(record[csvColumns["sortCategory"]]),
					StackSize:          csv.SafeStringToInt(record[csvColumns["stackSize"]]),
					BuyFromVendorPrice: adjustedBuyPrice,
					SellToVendorPrice:  csv.SafeStringToInt(record[csvColumns["sellToVendorPrice"]]),
					ClassJobCategory:   csv.SafeStringToInt(record[csvColumns["classJobCategory"]]),
					CanBeTraded:        canBeTraded,
					DropsFromDungeon:   csv.SafeStringToBool(record[csvColumns["dungeonDrop"]]),
					CanBeHq:            csv.SafeStringToBool(record[csvColumns["canBeHq"]]),
					CanDesynth:         canDesynth,
					IsCollectable:      csv.SafeStringToBool(record[csvColumns["alwaysCollectable"]]),
					IsGlamour:          csv.SafeStringToBool(record[csvColumns["isGlamourous"]]),
				}

				return &result, nil
			},
		},
	}
}

func isOldItem(itemId int, itemName string) bool {
	return itemId <= 1600 &&
		(strings.HasPrefix(itemName, "dated") ||
			strings.HasPrefix(itemName, "pair of dated"))
}
