package csv

import (
	"github.com/level-5-pidgey/MarketMoogleApi/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogleApi/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogleApi/csv/types"
)

func NewRecipeBookCsvReader() *csv.UngroupedXivApiCsvReader[csvType.RecipeBook] {
	csvColumns := map[string]int{
		"id":         0,
		"bookItemId": 1,
		"bookName":   2,
	}

	return &csv.UngroupedXivApiCsvReader[csvType.RecipeBook]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.RecipeBook]{
			FileName:   "SecretRecipeBook",
			RowsToSkip: 4,
			ProcessRow: func(record []string) (*csvType.RecipeBook, error) {
				return &csvType.RecipeBook{
					Id:         csv.SafeStringToInt(record[csvColumns["id"]]),
					BookItemId: csv.SafeStringToInt(record[csvColumns["bookItemId"]]),
					BookName:   record[csvColumns["bookName"]],
				}, nil
			},
		},
	}
}
