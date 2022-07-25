/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (recipeprovider.go) is part of MarketMoogleAPI and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package db

import (
	schema "MarketMoogleAPI/core/graph/model"
	"context"
)

type RecipeProvider interface {
	InsertRecipe(ctx context.Context, recipeToAdd *schema.Recipe) (*schema.Recipe, error)
	FindRecipesByItemId(ctx context.Context, itemId int) ([]*schema.Recipe, error)
	FindRecipeByObjectId(ctx context.Context, objectId string) (*schema.Recipe, error)
	GetAllRecipes(ctx context.Context) ([]*schema.Recipe, error)
}
