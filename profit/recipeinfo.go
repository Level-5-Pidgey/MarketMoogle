package profitCalc

import (
	"github.com/level-5-pidgey/MarketMoogle/csv/datacollection"
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
)

type RecipeInfo struct {
	Yield                  int
	CraftType              string
	RecipeLevel            int
	SpecializationRequired bool
	IsExpert               bool
	CanQuickSynth          bool
	SecretRecipeBook       string
	Craftsmanship          int
	Control                int
	Difficulty             int
	Durability             int
	Quality                int
	RecipeIngredients      []RecipeIngredients
}

type RecipeIngredients struct {
	ItemId   int
	Quantity int
}

func getRecipes(item *readertype.Item, collection *datacollection.DataCollection) (*[]RecipeInfo, error) {
	recipeLevels := *collection.RecipeLevels
	recipeBooks := *collection.RecipeBooks
	craftTypes := *collection.CraftTypes
	recipeMap := *collection.Recipes

	recipes, ok := recipeMap[item.Id]
	if !ok {
		// If no recipes are found for this item, that's okay, just return nothing
		return nil, nil
	}

	result := make([]RecipeInfo, 0, len(recipes))
	for _, recipe := range recipes {
		recipeBookRequired := ""
		recipeLevel := 1
		difficulty := 100
		durability := 60
		quality := 0
		craftType := ""

		if level, ok := recipeLevels[recipe.RecipeLevelId]; ok {
			recipeLevel = level.ClassJobLevel
			difficulty = level.Difficulty
			durability = level.Durability
			quality = level.Quality
		}

		if recipeBook, ok := recipeBooks[recipe.SecretRecipeBookId]; ok {
			recipeBookRequired = recipeBook.BookName
		}

		if recipeCraftType, ok := craftTypes[recipe.CraftType]; ok {
			craftType = recipeCraftType.Name
		}

		recipeInfo := RecipeInfo{
			Yield:                  recipe.Quantity,
			CraftType:              craftType,
			RecipeLevel:            recipeLevel,
			SpecializationRequired: recipe.SpecializationRequired,
			IsExpert:               recipe.IsExpert,
			CanQuickSynth:          recipe.CanQuickSynth,
			SecretRecipeBook:       recipeBookRequired,
			Craftsmanship:          recipe.RequiredCraftsmanship,
			Control:                recipe.RequiredControl,
			Difficulty:             difficulty,
			Durability:             durability,
			Quality:                quality,
			RecipeIngredients:      getRecipeIngredientsForItem(recipe),
		}

		result = append(result, recipeInfo)
	}

	return &result, nil
}

func getRecipeIngredientsForItem(recipe *readertype.Recipe) []RecipeIngredients {
	result := make([]RecipeIngredients, len(recipe.Ingredients))

	for i, recipeIngredient := range recipe.Ingredients {
		result[i] = RecipeIngredients{
			ItemId:   recipeIngredient.ItemId,
			Quantity: recipeIngredient.Quantity,
		}
	}

	return result
}
