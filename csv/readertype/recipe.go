package readertype

import (
	"github.com/level-5-pidgey/MarketMoogle/util"
)

type Ingredient struct {
	ItemId   int
	Quantity int
}

type Recipe struct {
	Id                     int
	CraftType              int
	RecipeLevelId          int
	ResultItemId           int
	Quantity               int
	Ingredients            []Ingredient
	RequiredCraftsmanship  int
	RequiredControl        int
	SpecializationRequired bool
	IsExpert               bool
	CanQuickSynth          bool
	SecretRecipeBookId     int
}

func (r Recipe) CreateFromCsvRow(record []string) (*Recipe, error) {
	ingredients := make([]Ingredient, 0)

	/*
		ingredients start at row 6 and are every 2 rows,
		quantities start at row 7 and are every 2 rows
	*/
	for i := 6; i < 24; i += 2 {
		ingredientItem := util.SafeStringToInt(record[i])
		quantity := util.SafeStringToInt(record[i+1])

		if ingredientItem < 1 || quantity < 1 {
			continue
		}

		ingredients = append(
			ingredients, Ingredient{
				ItemId:   ingredientItem,
				Quantity: quantity,
			},
		)
	}

	return &Recipe{
		Id:                     util.SafeStringToInt(record[0]),
		CraftType:              util.SafeStringToInt(record[2]),
		RecipeLevelId:          util.SafeStringToInt(record[3]),
		ResultItemId:           util.SafeStringToInt(record[4]),
		Quantity:               util.SafeStringToInt(record[5]),
		Ingredients:            ingredients,
		RequiredCraftsmanship:  util.SafeStringToInt(record[33]),
		RequiredControl:        util.SafeStringToInt(record[34]),
		SpecializationRequired: util.SafeStringToBool(record[44]),
		IsExpert:               util.SafeStringToBool(record[45]),
		CanQuickSynth:          util.SafeStringToBool(record[39]),
		SecretRecipeBookId:     util.SafeStringToInt(record[37]),
	}, nil
}

func (r Recipe) GetKey() int {
	return r.ResultItemId
}
