package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	generated "MarketMoogleAPI/core/graph/gen"
	schema "MarketMoogleAPI/core/graph/model"
	"MarketMoogleAPI/infrastructure/providers/db"
	"context"
	"fmt"
)

// Users is the resolver for the Users field.
func (r *queryResolver) Users(ctx context.Context) ([]*schema.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Items is the resolver for the Items field.
func (r *queryResolver) Items(ctx context.Context) ([]*schema.Item, error) {
	return r.itemProv.GetAllItems(ctx)
}

// MarketboardEntries is the resolver for the MarketboardEntries field.
func (r *queryResolver) MarketboardEntries(ctx context.Context) ([]*schema.MarketboardEntry, error) {
	return r.mbProv.GetAllMarketboardEntries(ctx)
}

// Recipes is the resolver for the Recipes field.
func (r *queryResolver) Recipes(ctx context.Context) ([]*schema.Recipe, error) {
	return r.recipeProv.GetAllRecipes(ctx)
}

// GetUser is the resolver for the GetUser field.
func (r *queryResolver) GetUser(ctx context.Context, userID int) (*schema.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// GetItem is the resolver for the GetItem field.
func (r *queryResolver) GetItem(ctx context.Context, itemID int) (*schema.Item, error) {
	return r.itemProv.FindItemByItemId(ctx, itemID)
}

// GetMarketboardEntriesForItem is the resolver for the GetMarketboardEntriesForItem field.
func (r *queryResolver) GetMarketboardEntriesForItem(ctx context.Context, itemID int) ([]*schema.MarketboardEntry, error) {
	return r.mbProv.GetAllMarketboardEntries(ctx)
}

// GetRecipesForItem is the resolver for the GetRecipesForItem field.
func (r *queryResolver) GetRecipesForItem(ctx context.Context, itemID int) ([]*schema.Recipe, error) {
	return r.recipeProv.FindRecipesByItemId(ctx, itemID)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{
		Resolver:   r,
		recipeProv: db.NewRecipeDatabaseProvider(r.DbClient),
		mbProv:     db.NewMarketboardDatabaseProvider(r.DbClient),
		itemProv:   db.NewItemDataBaseProvider(r.DbClient),
	}
}

type queryResolver struct {
	*Resolver
	recipeProv *db.RecipeDatabaseProvider
	mbProv     *db.MarketboardDatabaseProvider
	itemProv   *db.ItemDatabaseProvider
}
