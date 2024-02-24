package datacollection

import (
	"github.com/level-5-pidgey/MarketMoogle/domain"
)

type RecipeDataCollection struct {
	Recipes      *map[int][]domain.Recipe
	RecipeBooks  *map[int]domain.RecipeBook
	RecipeLevels *map[int]domain.RecipeLevel
	CraftTypes   *map[int]domain.CraftType
}
