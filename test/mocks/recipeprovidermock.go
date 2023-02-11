/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (recipeprovidermock.go) is part of MarketMoogleAPI and is released under the GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package mocks

import (
	schema "MarketMoogleAPI/core/graph/model"
	"context"
)

type TestRecipeProvider struct {
	RecipeDatabase map[int]*schema.Recipe
}

func (recipeProv TestRecipeProvider) InsertRecipe(ctx context.Context, recipeToAdd *schema.Recipe) (*schema.Recipe, error) {
	recipeProv.RecipeDatabase[recipeToAdd.ItemResultID] = recipeToAdd

	return recipeToAdd, nil
}

func (recipeProv TestRecipeProvider) FindRecipesByItemId(ctx context.Context, itemId int) ([]*schema.Recipe, error) {
	if val, ok := recipeProv.RecipeDatabase[itemId]; ok {
		return []*schema.Recipe{val}, nil
	}

	return []*schema.Recipe{}, nil
}

func (recipeProv TestRecipeProvider) FindRecipeByObjectId(ctx context.Context, objectId string) (*schema.Recipe, error) {
	//ObjectID isn't on the schema objects, so we cannot search by them without actually querying mongo
	//This shouldn't be an issue since this method isn't used for much
	return nil, nil
}

func (recipeProv TestRecipeProvider) GetAllRecipes(ctx context.Context) ([]*schema.Recipe, error) {
	var results []*schema.Recipe
	for _, recipe := range recipeProv.RecipeDatabase {
		results = append(results, recipe)
	}

	return results, nil
}
