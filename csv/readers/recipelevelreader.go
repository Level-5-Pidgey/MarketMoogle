package csv

import (
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/csv/types"
)

func NewRecipeLevelReader() *csv.UngroupedXivApiCsvReader[csvType.RecipeLevel] {
	csvColumns := map[string]int{
		"key":                    0,
		"classJobLevel":          1,
		"stars":                  2,
		"suggestedCraftsmanship": 3,
		"suggestedControl":       4,
		"difficulty":             5,
		"quality":                6,
		"durability":             11,
	}

	return &csv.UngroupedXivApiCsvReader[csvType.RecipeLevel]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.RecipeLevel]{
			FileName:   "RecipeLevelTable",
			RowsToSkip: 4,
			ProcessRow: func(record []string) (*csvType.RecipeLevel, error) {
				return &csvType.RecipeLevel{
					Id:                     csv.SafeStringToInt(record[csvColumns["key"]]),
					ClassJobLevel:          csv.SafeStringToInt(record[csvColumns["classJobLevel"]]),
					SuggestedCraftsmanship: csv.SafeStringToInt(record[csvColumns["suggestedCraftsmanship"]]),
					SuggestedControl:       csv.SafeStringToInt(record[csvColumns["suggestedControl"]]),
					Difficulty:             csv.SafeStringToInt(record[csvColumns["difficulty"]]),
					Quality:                csv.SafeStringToInt(record[csvColumns["quality"]]),
					Durability:             csv.SafeStringToInt(record[csvColumns["durability"]]),
				}, nil
			},
		},
	}
}
