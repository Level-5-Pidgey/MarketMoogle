/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (query.resolvers.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	generated "MarketMoogleAPI/core/graph/gen"
	schema "MarketMoogleAPI/core/graph/model"
	"context"
	"fmt"
)

func (r *queryResolver) Users(ctx context.Context) ([]*schema.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Items(ctx context.Context) ([]*schema.Item, error) {
	return dbProv.GetAllItems(ctx)
}

func (r *queryResolver) MarketboardEntries(ctx context.Context) ([]*schema.MarketboardEntry, error) {
	return dbProv.GetAllMarketboardEntries(ctx)
}

func (r *queryResolver) Recipes(ctx context.Context) ([]*schema.Recipe, error) {
	return dbProv.GetAllRecipes(ctx)
}

func (r *queryResolver) GetUser(ctx context.Context, userID int) (*schema.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetItem(ctx context.Context, itemID int) (*schema.Item, error) {
	return dbProv.FindItemByItemId(ctx, itemID)
}

func (r *queryResolver) GetMarketboardEntriesForItem(ctx context.Context, itemID int) ([]*schema.MarketboardEntry, error) {
	return dbProv.FindMarketboardEntriesByItemId(ctx, itemID)
}

func (r *queryResolver) GetRecipesForItem(ctx context.Context, itemID int) ([]*schema.Recipe, error) {
	return dbProv.FindRecipesByItemId(ctx, itemID)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
