/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (marketboardprovider.go) is part of MarketMoogleAPI and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package database

import (
	"MarketMoogle/core/apitypes/universalis"
	schema "MarketMoogle/core/graph/model"
	"context"
)

type MarketBoardProvider interface {
	CreateMarketEntry(ctx context.Context, entryFromApi *schema.MarketboardEntry) (*schema.MarketboardEntry, error)
	ReplaceMarketEntry(ctx context.Context, itemId int, dataCenter string, newEntry *universalis.MarketQuery, currentTimestamp *string) error
	FindMarketboardEntryByObjectId(ctx context.Context, objectId string) (*schema.MarketboardEntry, error)
	FindMarketboardEntriesByItemId(ctx context.Context, itemId int) ([]*schema.MarketboardEntry, error)
	FindItemEntryAcrossDataCenter(ctx context.Context, itemId int, dataCenter string) (*schema.MarketboardEntry, error)
	GetAllMarketboardEntries(ctx context.Context) ([]*schema.MarketboardEntry, error)
}
