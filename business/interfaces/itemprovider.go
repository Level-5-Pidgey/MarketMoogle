/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (itemprovider.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package interfaces

import (
	schema "MarketMoogleAPI/core/graph/model"
	"golang.org/x/net/context"
)

type ItemProvider interface {
	SaveItem(input *schema.Item) (*schema.Item, error)
	FindItemByObjectId(ID string) (*schema.Item, error)
	FindItemByItemId(ItemId int) (*schema.Item, error)
	GetAllItems(ctx context.Context) ([]*schema.Item, error)
	GetItemFromApi(itemID *int) (*schema.Item, error)
	SaveItemFromApi(itemID *int) (*schema.Item, error)
}
