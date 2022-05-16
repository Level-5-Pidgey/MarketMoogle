/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (marketboardprovider.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package interfaces

import (
	"MarketMoogleAPI/core/apitypes/universalis"
	schema "MarketMoogleAPI/core/graph/model"
	"context"
)

type MarketBoardProvider interface {
	CreateMarketboardEntryFromApi(dataCenter *string, itemId *int) (*schema.MarketboardEntry, error)
	ReplaceMarketEntries(itemId *int, dataCenter *string, newEntry *universalis.MarketQuery, currentTimestamp *string) error
	SaveMarketboardEntry(input *schema.MarketboardEntry) (*schema.MarketboardEntry, error)
	FindMarketboardEntryByObjectId(ID string) (*schema.MarketboardEntry, error)
	FindMarketboardEntriesByItemId(ctx context.Context, itemId int) ([]*schema.MarketboardEntry, error)
	GetAllMarketboardEntries(ctx context.Context) ([]*schema.MarketboardEntry, error)
	GetCrossDcResaleProfit(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string) (int, error)
}
