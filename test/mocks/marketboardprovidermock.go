/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (marketboardprovidermock.go) is part of MarketMoogleAPI and is released under the GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package mocks

import (
	"MarketMoogleAPI/core/apitypes/universalis"
	schema "MarketMoogleAPI/core/graph/model"
	"context"
	"errors"
)

type TestMarketboardProvider struct {
	MbDatabase  map[int]*schema.MarketboardEntry
	maxListings int
}

func (mbProv TestMarketboardProvider) CreateMarketEntry(ctx context.Context, entryFromApi *schema.MarketboardEntry) (*schema.MarketboardEntry, error) {
	mbProv.MbDatabase[entryFromApi.ItemID] = entryFromApi

	return entryFromApi, nil
}

func (mbProv TestMarketboardProvider) ReplaceMarketEntry(ctx context.Context, itemId int, dataCenter string, newEntry *universalis.MarketQuery, currentTimestamp *string) error {
	updatedEntry := schema.MarketboardEntry{
		ItemID:              itemId,
		LastUpdateTime:      *currentTimestamp,
		MarketEntries:       newEntry.GetMarketEntries(mbProv.maxListings),
		MarketHistory:       newEntry.GetItemHistory(mbProv.maxListings),
		DataCenter:          dataCenter,
		CurrentAveragePrice: newEntry.CurrentAveragePrice,
		CurrentMinPrice:     &newEntry.MinPrice,
		RegularSaleVelocity: newEntry.RegularSaleVelocity,
		HqSaleVelocity:      newEntry.HqSaleVelocity,
		NqSaleVelocity:      newEntry.NqSaleVelocity,
	}

	mbProv.MbDatabase[itemId] = &updatedEntry

	return nil
}

func (mbProv TestMarketboardProvider) FindMarketboardEntryByObjectId(ctx context.Context, objectId string) (*schema.MarketboardEntry, error) {
	//ObjectID isn't on the schema objects, so we cannot search by them without actually querying mongo
	//This shouldn't be an issue since this method isn't used for much
	return nil, nil
}

func (mbProv TestMarketboardProvider) FindMarketboardEntriesByItemId(ctx context.Context, itemId int) ([]*schema.MarketboardEntry, error) {
	if val, ok := mbProv.MbDatabase[itemId]; ok {
		return []*schema.MarketboardEntry{val}, nil
	}

	return nil, errors.New("unable to find market board entries for items")
}

func (mbProv TestMarketboardProvider) FindItemEntryOnDc(ctx context.Context, itemId int, dataCenter string) (*schema.MarketboardEntry, error) {
	if val, ok := mbProv.MbDatabase[itemId]; ok {
		if val.DataCenter == dataCenter {
			return val, nil
		} else {
			return nil, errors.New("unable to find entry on datacenter for given item")
		}
	}

	//Return nil if nothing has been found
	return nil, nil
}

func (mbProv TestMarketboardProvider) GetAllMarketboardEntries(ctx context.Context) ([]*schema.MarketboardEntry, error) {
	var results []*schema.MarketboardEntry
	for _, item := range mbProv.MbDatabase {
		results = append(results, item)
	}

	return results, nil
}
