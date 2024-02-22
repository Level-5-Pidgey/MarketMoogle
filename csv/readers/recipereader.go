package csv

import (
	"fmt"
	"github.com/level-5-pidgey/MarketMoogleApi/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogleApi/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogleApi/csv/types"
)

func NewRecipeCsvReader() *csv.GroupedXivApiCsvReader[csvType.Recipe] {
	csvColumns := map[string]int{
		"id":                     0,
		"craftType":              2,
		"level":                  3,
		"resultItem":             4,
		"resultAmount":           5,
		"ingredient0":            6,
		"quantity0":              7,
		"ingredient1":            8,
		"quantity1":              9,
		"ingredient2":            10,
		"quantity2":              11,
		"ingredient3":            12,
		"quantity3":              13,
		"ingredient4":            14,
		"quantity4":              15,
		"ingredient5":            16,
		"quantity5":              17,
		"ingredient6":            18,
		"quantity6":              19,
		"ingredient7":            20,
		"quantity7":              21,
		"ingredient8":            22,
		"quantity8":              23,
		"ingredient9":            24,
		"quantity9":              25,
		"recipeBookId":           37,
		"requiredCraftsmanship":  33,
		"requiredControl":        34,
		"canQuickSynth":          39,
		"specializationRequired": 44,
		"isExpert":               45,
	}

	return &csv.GroupedXivApiCsvReader[csvType.Recipe]{
		XivApiCsvInfo: csvInterface.XivApiCsvInfo[csvType.Recipe]{
			FileName:   "Recipe",
			RowsToSkip: 4,
			ProcessRow: func(record []string) (*csvType.Recipe, error) {
				ingredients := make([]csvType.Ingredient, 0)

				for i := 0; i < 10; i++ {
					ingredientItem := csv.SafeStringToInt(record[csvColumns[fmt.Sprintf("ingredient%d", i)]])
					quantity := csv.SafeStringToInt(record[csvColumns[fmt.Sprintf("quantity%d", i)]])

					if ingredientItem < 1 || quantity < 1 {
						continue
					}

					ingredients = append(
						ingredients, csvType.Ingredient{
							ItemId:   ingredientItem,
							Quantity: quantity,
						},
					)
				}

				return &csvType.Recipe{
					Id:                     csv.SafeStringToInt(record[csvColumns["id"]]),
					CraftType:              csv.SafeStringToInt(record[csvColumns["craftType"]]),
					RecipeLevelId:          csv.SafeStringToInt(record[csvColumns["level"]]),
					ResultItemId:           csv.SafeStringToInt(record[csvColumns["resultItem"]]),
					Quantity:               csv.SafeStringToInt(record[csvColumns["resultAmount"]]),
					Ingredients:            ingredients,
					RequiredCraftsmanship:  csv.SafeStringToInt(record[csvColumns["requiredCraftsmanship"]]),
					RequiredControl:        csv.SafeStringToInt(record[csvColumns["requiredControl"]]),
					SpecializationRequired: csv.SafeStringToBool(record[csvColumns["specializationRequired"]]),
					IsExpert:               csv.SafeStringToBool(record[csvColumns["isExpert"]]),
					CanQuickSynth:          csv.SafeStringToBool(record[csvColumns["canQuickSynth"]]),
					SecretRecipeBookId:     csv.SafeStringToInt(record[csvColumns["recipeBookId"]]),
				}, nil
			},
		},
	}
}
