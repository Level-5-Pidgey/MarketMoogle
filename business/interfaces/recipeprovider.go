/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (recipeprovider.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package interfaces

import (
	schema "MarketMoogleAPI/core/graph/model"
	"context"
)

type RecipeProvider interface {
	CreateRecipesFromApi(itemID *int) (*[]*schema.Recipe, error)
	SaveRecipes(input *[]*schema.Recipe) (*[]*schema.Recipe, error)
	FindRecipesByItemId(ctx context.Context, ItemId int) ([]*schema.Recipe, error)
	FindRecipeByObjectId(ID string) (*schema.Recipe, error)
	GetAllRecipes(ctx context.Context) ([]*schema.Recipe, error)
}
