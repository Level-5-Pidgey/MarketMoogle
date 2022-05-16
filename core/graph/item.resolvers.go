/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (item.resolvers.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	generated "MarketMoogleAPI/core/graph/gen"
	schema "MarketMoogleAPI/core/graph/model"
	"context"
)

func (r *itemResolver) Recipes(ctx context.Context, obj *schema.Item) ([]*schema.Recipe, error) {
	return dbProv.FindRecipesByItemId(ctx, obj.ItemID)
}

func (r *itemResolver) MarketboardEntries(ctx context.Context, obj *schema.Item) ([]*schema.MarketboardEntry, error) {
	return dbProv.FindMarketboardEntriesByItemId(ctx, obj.ItemID)
}

func (r *itemResolver) ResaleValue(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string) (int, error) {
	return dbProv.GetCrossDcResaleProfit(ctx, obj, &dataCenter, &homeServer)
}

func (r *itemResolver) VendorResaleValue(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string) (int, error) {
	return dbProv.GetVendorResaleProfit(ctx, obj, &dataCenter, &homeServer)
}

func (r *itemResolver) RecipeResaleValue(ctx context.Context, obj *schema.Item, buyFromOtherSevers *bool, buyCrystals *bool, dataCenter string, homeServer string) (*schema.RecipeResaleInformation, error) {
	return dbProv.GetCraftingProfit(ctx, obj, &dataCenter, &homeServer, buyCrystals, buyFromOtherSevers)
}

// Item returns generated.ItemResolver implementation.
func (r *Resolver) Item() generated.ItemResolver { return &itemResolver{r} }

type itemResolver struct{ *Resolver }
