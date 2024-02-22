package csv

type RecipeLevel struct {
	Id                     int
	ClassJobLevel          int
	SuggestedCraftsmanship int
	SuggestedControl       int
	Difficulty             int
	Quality                int
	Durability             int
}

func (r RecipeLevel) GetKey() int {
	return r.Id
}
