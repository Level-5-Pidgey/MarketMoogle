/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (itemprovider.go) is part of MarketMoogleAPI and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package database

import (
	schema "MarketMoogleAPI/core/graph/model"
	"golang.org/x/net/context"
)

type ItemProvider interface {
	InsertItem(ctx context.Context, input *schema.Item) (*schema.Item, error)
	FindItemByObjectId(ctx context.Context, ID string) (*schema.Item, error)
	FindItemByItemId(ctx context.Context, itemId int) (*schema.Item, error)
	GetAllItems(ctx context.Context) ([]*schema.Item, error)
	UpdateVendorSellPrice(ctx context.Context, itemId int, newPrice int) error
	UpdateLevequestInfo(ctx context.Context, itemId int, leve schema.CraftingLeve) error
}
