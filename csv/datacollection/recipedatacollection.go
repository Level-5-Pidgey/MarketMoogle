package datacollection

import csvType "github.com/level-5-pidgey/MarketMoogle/csv/types"

type RecipeDataCollection struct {
	Recipes      *map[int][]csvType.Recipe
	RecipeBooks  *map[int]csvType.RecipeBook
	RecipeLevels *map[int]csvType.RecipeLevel
	CraftTypes   *map[int]csvType.CraftType
}
