/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (recipecontents.resolvers.go) is part of MarketMoogle and is released GNU General Public License.
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

func (r *recipeContentsResolver) Item(ctx context.Context, obj *schema.RecipeContents) (*schema.Item, error) {
	item, err := dbProv.FindItemByItemId(ctx, obj.ItemID)

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
	return &recipeContentsResolver{r}
}

type recipeContentsResolver struct{ *Resolver }
