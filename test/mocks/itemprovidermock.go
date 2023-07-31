/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (itemprovidermock.go) is part of MarketMoogleAPI and is released under the GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package mocks

import (
	schema "MarketMoogleAPI/core/graph/model"
	"context"
	"errors"
)

type TestItemProvider struct {
	ItemDatabase map[int]*schema.Item
}

func (itemProv TestItemProvider) InsertItem(ctx context.Context, input *schema.Item) (*schema.Item, error) {
	if input == nil {
		return input, nil
	}
	
	itemProv.ItemDatabase[input.Id] = input

	return input, nil
}

func (itemProv TestItemProvider) FindItemByObjectId(ctx context.Context, ID string) (*schema.Item, error) {
	//ObjectID isn't on the schema objects, so we cannot search by them without actually querying mongo
	//This shouldn't be an issue since this method isn't used for much
	return nil, nil
}

func (itemProv TestItemProvider) FindItemByItemId(ctx context.Context, itemId int) (*schema.Item, error) {
	if val, ok := itemProv.ItemDatabase[itemId]; ok {
		return val, nil
	}

	return nil, nil
}

func (itemProv TestItemProvider) GetAllItems(ctx context.Context) ([]*schema.Item, error) {
	var results []*schema.Item
	for _, item := range itemProv.ItemDatabase {
		results = append(results, item)
	}

	return results, nil
}

func (itemProv TestItemProvider) UpdateVendorSellPrice(ctx context.Context, itemId int, newPrice int) error {
	if val, ok := itemProv.ItemDatabase[itemId]; ok {
		val.BuyFromVendorValue = &newPrice

		return nil
	}

	return errors.New("unable to find item in test database")
}

func (itemProv TestItemProvider) UpdateLevequestInfo(ctx context.Context, itemId int, leve schema.CraftingLeve) error {
	if val, ok := itemProv.ItemDatabase[itemId]; ok {
		val.AssociatedLeve = &leve

		return nil
	}

	return errors.New("unable to find item in test database")
}
