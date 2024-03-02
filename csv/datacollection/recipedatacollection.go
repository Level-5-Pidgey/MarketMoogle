package datacollection

import "github.com/level-5-pidgey/MarketMoogle/csv/readertype"

type RecipeDataCollection struct {
	Recipes      *map[int][]*readertype.Recipe
	RecipeBooks  *map[int]*readertype.RecipeBook
	RecipeLevels *map[int]*readertype.RecipeLevel
	CraftTypes   *map[int]*readertype.CraftType
}
