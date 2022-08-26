/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (profitprovider.go) is part of MarketMoogleAPI and is released under the GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package Business

import (
	schema "MarketMoogleAPI/core/graph/model"
	"context"
)

type ProfitCalculator interface {
	GetVendorFlipProfit(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string) (int, error)
	GetResaleInfoForItem(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string, buyCrystals *bool, buyFromOtherServers *bool) (*schema.RecipeResaleInfo, error)
	GetCrossDcResaleProfit(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string) (*schema.ResaleInfo, error)
	GetCheapestOnServer(entry *schema.MarketboardEntry, server string) *schema.MarketEntry
	GetRecipePurchaseInfo(componentItem *schema.Item, mbEntry *schema.MarketboardEntry, homeServer string, buyFromOtherServers *bool) *schema.RecipePurchaseInfo
	GetCheapestOnDc(entry *schema.MarketboardEntry) *schema.MarketEntry
}
