package domain

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

type Ingredient struct {
	ItemId   int
	Quantity int
}

func (r Recipe) GetKey() int {
	return r.ResultItemId
}
