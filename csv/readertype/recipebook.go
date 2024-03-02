package readertype

import "github.com/level-5-pidgey/MarketMoogle/util"

type RecipeBook struct {
	Id         int
	BookItemId int
	BookName   string
}

func (r RecipeBook) CreateFromCsvRow(record []string) (*RecipeBook, error) {
	return &RecipeBook{
		Id:         util.SafeStringToInt(record[0]),
		BookItemId: util.SafeStringToInt(record[1]),
		BookName:   record[2],
	}, nil
}

func (r RecipeBook) GetKey() int {
	return r.Id
}
