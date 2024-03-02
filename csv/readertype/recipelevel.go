package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type RecipeLevel struct {
	Id                     int
	ClassJobLevel          int
	SuggestedCraftsmanship int
	SuggestedControl       int
	Difficulty             int
	Quality                int
	Durability             int
}

func (r RecipeLevel) CreateFromCsvRow(record []string) (*RecipeLevel, error) {
	return &RecipeLevel{
		Id:                     util.SafeStringToInt(record[0]),
		ClassJobLevel:          util.SafeStringToInt(record[1]),
		SuggestedCraftsmanship: util.SafeStringToInt(record[3]),
		SuggestedControl:       util.SafeStringToInt(record[4]),
		Difficulty:             util.SafeStringToInt(record[5]),
		Quality:                util.SafeStringToInt(record[6]),
		Durability:             util.SafeStringToInt(record[11]),
	}, nil
}

func (r RecipeLevel) GetKey() int {
	return r.Id
}
