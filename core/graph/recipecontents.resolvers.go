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

// Item is the resolver for the Item field.
func (r *recipeContentsResolver) Item(ctx context.Context, obj *schema.RecipeContents) (*schema.Item, error) {
	item, err := r.itemProv.FindItemByItemId(ctx, obj.ItemID)

	if item == nil {
		dummyDesc := fmt.Sprintf("Item with ID %d could not be found.", obj.ItemID)
		dummyItem := schema.Item{
			ItemID:             0,
			Name:               "Item not found",
			Description:        &dummyDesc,
			CanBeHq:            false,
			IconID:             0,
			SellToVendorValue:  nil,
			BuyFromVendorValue: nil,
		}

		return &dummyItem, nil
	}

	return item, err
}

// RecipeContents returns generated.RecipeContentsResolver implementation.
func (r *Resolver) RecipeContents() generated.RecipeContentsResolver {
	iProv := db.NewItemDataBaseProvider(r.DbClient)
	
	return &recipeContentsResolver{
		Resolver: 	r,
		itemProv:   iProv,
	}
}

type recipeContentsResolver struct{ 
	*Resolver
	itemProv   *db.ItemDatabaseProvider
}
