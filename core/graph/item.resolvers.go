package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	generated "MarketMoogleAPI/core/graph/gen"
	schema "MarketMoogleAPI/core/graph/model"
	"MarketMoogleAPI/infrastructure/providers"
	"MarketMoogleAPI/infrastructure/providers/db"
	"context"
)

// Recipes is the resolver for the Recipes field.
func (r *itemResolver) Recipes(ctx context.Context, obj *schema.Item) ([]*schema.Recipe, error) {
	return r.recipeProv.FindRecipesByItemId(ctx, obj.ItemID)
}

// MarketboardEntries is the resolver for the MarketboardEntries field.
func (r *itemResolver) MarketboardEntries(ctx context.Context, obj *schema.Item) ([]*schema.MarketboardEntry, error) {
	return r.mbProv.FindMarketboardEntriesByItemId(ctx, obj.ItemID)
}

// ResaleValue is the resolver for the ResaleValue field.
func (r *itemResolver) ResaleValue(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string) (int, error) {
	return r.profitProv.GetCrossDcResaleProfit(ctx, obj, dataCenter, homeServer)
}

// VendorResaleValue is the resolver for the VendorResaleValue field.
func (r *itemResolver) VendorResaleValue(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string) (int, error) {
	return r.profitProv.GetVendorResaleProfit(ctx, obj, dataCenter, homeServer)
}

// RecipeResaleValue is the resolver for the RecipeResaleValue field.
func (r *itemResolver) RecipeResaleValue(ctx context.Context, obj *schema.Item, buyFromOtherSevers *bool, buyCrystals *bool, dataCenter string, homeServer string) (*schema.RecipeResaleInformation, error) {
	return r.profitProv.GetCraftingProfit(ctx, obj, dataCenter, homeServer, buyCrystals, buyFromOtherSevers)
}

// Item returns generated.ItemResolver implementation.
func (r *Resolver) Item() generated.ItemResolver {

	rProv := db.NewRecipeDatabaseProvider(r.DbClient)
	mProv := db.NewMarketboardDatabaseProvider(r.DbClient)
	iProv := db.NewItemDataBaseProvider(r.DbClient)

	return &itemResolver{
		Resolver:   r,
		recipeProv: rProv,
		mbProv:     mProv,
		itemProv:   iProv,
		profitProv: providers.NewItemProfitProvider(rProv, mProv, iProv),
	}
}

type itemResolver struct {
	*Resolver
	recipeProv *db.RecipeDatabaseProvider
	mbProv     *db.MarketboardDatabaseProvider
	itemProv   *db.ItemDatabaseProvider
	profitProv *providers.ItemProfitProvider
}
