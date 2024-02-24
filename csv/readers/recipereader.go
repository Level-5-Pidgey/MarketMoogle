package csv

import (
	"fmt"
	"github.com/level-5-pidgey/MarketMoogle/csv"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvType "github.com/level-5-pidgey/MarketMoogle/domain"
	"github.com/level-5-pidgey/MarketMoogle/util"
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
					ingredientItem := util.SafeStringToInt(record[csvColumns[fmt.Sprintf("ingredient%d", i)]])
					quantity := util.SafeStringToInt(record[csvColumns[fmt.Sprintf("quantity%d", i)]])

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
					Id:                     util.SafeStringToInt(record[csvColumns["id"]]),
					CraftType:              util.SafeStringToInt(record[csvColumns["craftType"]]),
					RecipeLevelId:          util.SafeStringToInt(record[csvColumns["level"]]),
					ResultItemId:           util.SafeStringToInt(record[csvColumns["resultItem"]]),
					Quantity:               util.SafeStringToInt(record[csvColumns["resultAmount"]]),
					Ingredients:            ingredients,
					RequiredCraftsmanship:  util.SafeStringToInt(record[csvColumns["requiredCraftsmanship"]]),
					RequiredControl:        util.SafeStringToInt(record[csvColumns["requiredControl"]]),
					SpecializationRequired: util.SafeStringToBool(record[csvColumns["specializationRequired"]]),
					IsExpert:               util.SafeStringToBool(record[csvColumns["isExpert"]]),
					CanQuickSynth:          util.SafeStringToBool(record[csvColumns["canQuickSynth"]]),
					SecretRecipeBookId:     util.SafeStringToInt(record[csvColumns["recipeBookId"]]),
				}, nil
			},
		},
	}
}
