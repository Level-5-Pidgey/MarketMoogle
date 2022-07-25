/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (profitprovider.go) is part of MarketMoogleAPI and is released under the GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package business

import (
	schema "MarketMoogleAPI/core/graph/model"
	"context"
)

type ProfitCalculator interface {
	GetItemValue(componentItem *schema.Item, mbEntry *schema.MarketboardEntry, homeServer string, buyFromOtherServers *bool) (int, string)
	GetVendorResaleProfit(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string) (int, error)
	GetCraftingProfit(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string, buyCrystals *bool, buyFromOtherServers *bool) (*schema.RecipeResaleInformation, error)
	GetCrossDcResaleProfit(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string) (int, error)
	GetCheapestOnServer(entry *schema.MarketboardEntry, server string) int
	GetCheapestPriceAndServer(entry *schema.MarketboardEntry) (int, string)
}
