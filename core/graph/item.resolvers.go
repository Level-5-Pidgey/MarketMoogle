package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	generated "MarketMoogleAPI/core/graph/gen"
	schema "MarketMoogleAPI/core/graph/model"
	"MarketMoogleAPI/infrastructure/providers"
	"MarketMoogleAPI/infrastructure/providers/database"
	"context"
	"log"
)

// Recipes is the resolver for the Recipes field.
func (r *itemResolver) Recipes(ctx context.Context, obj *schema.Item) ([]*schema.Recipe, error) {
	return r.recipeProv.FindRecipesByItemId(ctx, obj.Id)
}

// MarketboardEntries is the resolver for the MarketboardEntries field.
func (r *itemResolver) MarketboardEntries(ctx context.Context, obj *schema.Item) ([]*schema.MarketboardEntry, error) {
	return r.mbProv.FindMarketboardEntriesByItemId(ctx, obj.Id)
}

// LeveProfit is the resolver for the LeveProfit field.
func (r *itemResolver) LeveProfit(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string, buyFromOtherSevers *bool) (*schema.ProfitInfo, error) {
	panic("implement me")
}

// DcFlipProfit is the resolver for the DcFlipProfit field.
func (r *itemResolver) DcFlipProfit(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string) (*schema.ProfitInfo, error) {
	resaleProfit, err := r.profitProv.GetCrossDcResaleProfit(ctx, obj, dataCenter, homeServer)

	if err != nil {
		log.Fatal(err)
	}

	return resaleProfit, nil
}

// VendorFlipProfit is the resolver for the VendorFlipProfit field.
func (r *itemResolver) VendorFlipProfit(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string) (*schema.ProfitInfo, error) {
	vendorPrice := 0
	if obj.BuyFromVendorValue != nil {
		vendorPrice = *obj.BuyFromVendorValue
	}

	purchaseInfo := schema.ItemCostInfo{
		Item:            obj,
		ServerToBuyFrom: homeServer,
		BuyFromVendor:   true,
		PricePer:        vendorPrice,
		TotalCost:       vendorPrice,
		Quantity:        1,
	}

	resaleProfit, err := r.profitProv.GetVendorFlipProfit(ctx, obj, dataCenter, homeServer)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &schema.ProfitInfo{
		Profit:          resaleProfit,
		ItemID:          obj.Id,
		Quantity:        1,
		SingleCost:      resaleProfit,
		TotalCost:       resaleProfit,
		ItemsToPurchase: []*schema.ItemCostInfo{&purchaseInfo},
	}, nil
}

// RecipeProfit is the resolver for the RecipeProfit field.
func (r *itemResolver) RecipeProfit(ctx context.Context, obj *schema.Item, buyFromOtherSevers *bool, buyCrystals *bool, dataCenter string, homeServer string) (*schema.RecipeProfitInfo, error) {
	recipeResaleInfo, err := r.profitProv.GetRecipeProfitForItem(ctx, obj, dataCenter, homeServer, buyCrystals, buyFromOtherSevers)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return recipeResaleInfo, nil
}

func (r *Resolver) Item() generated.ItemResolver {
	rProv := database.NewRecipeDatabaseProvider(r.DbClient)
	iProv := database.NewItemDataBaseProvider(r.DbClient)
	mProv := database.NewMarketboardDatabaseProvider(r.DbClient)

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
	recipeProv *database.RecipeDatabaseProvider
	mbProv     *database.MarketboardDatabaseProvider
	itemProv   *database.ItemDatabaseProvider
	profitProv *providers.ItemProfitProvider
}
