/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (profitprovider_test.go) is part of MarketMoogleAPI and is released under the GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package test

import (
	schema "MarketMoogleAPI/core/graph/model"
	"MarketMoogleAPI/infrastructure/providers"
	"MarketMoogleAPI/infrastructure/providers/database"
	"context"
	"reflect"
	"testing"
)

func TestItemProfitProvider_GetCheapestOnDc(t *testing.T) {
	type fields struct {
		maxValue            int
		returnUnlisted      bool
		recipeProvider      *database.RecipeDatabaseProvider
		marketboardProvider *database.MarketboardDatabaseProvider
		itemProvider        *database.ItemDatabaseProvider
	}

	type args struct {
		entry *schema.MarketboardEntry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *schema.MarketEntry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profitProv := providers.ItemProfitProvider{
				maxValue:            tt.fields.maxValue,
				returnUnlisted:      tt.fields.returnUnlisted,
				recipeProvider:      tt.fields.recipeProvider,
				marketboardProvider: tt.fields.marketboardProvider,
				itemProvider:        tt.fields.itemProvider,
			}
			if got := profitProv.GetCheapestOnDc(tt.args.entry); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCheapestOnDc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemProfitProvider_GetCheapestOnServer(t *testing.T) {
	type fields struct {
		maxValue            int
		returnUnlisted      bool
		recipeProvider      *database.RecipeDatabaseProvider
		marketboardProvider *database.MarketboardDatabaseProvider
		itemProvider        *database.ItemDatabaseProvider
	}
	type args struct {
		entry  *schema.MarketboardEntry
		server string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *schema.MarketEntry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profitProv := providers.ItemProfitProvider{
				maxValue:            tt.fields.maxValue,
				returnUnlisted:      tt.fields.returnUnlisted,
				recipeProvider:      tt.fields.recipeProvider,
				marketboardProvider: tt.fields.marketboardProvider,
				itemProvider:        tt.fields.itemProvider,
			}

			if got := profitProv.GetCheapestOnServer(tt.args.entry, tt.args.server); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCheapestOnServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemProfitProvider_GetCraftingProfit(t *testing.T) {
	type fields struct {
		maxValue            int
		returnUnlisted      bool
		recipeProvider      *database.RecipeDatabaseProvider
		marketboardProvider *database.MarketboardDatabaseProvider
		itemProvider        *database.ItemDatabaseProvider
	}
	type args struct {
		ctx                 context.Context
		obj                 *schema.Item
		dataCenter          string
		homeServer          string
		buyCrystals         *bool
		buyFromOtherServers *bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *schema.RecipeResaleInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profitProv := providers.ItemProfitProvider{
				maxValue:            tt.fields.maxValue,
				returnUnlisted:      tt.fields.returnUnlisted,
				recipeProvider:      tt.fields.recipeProvider,
				marketboardProvider: tt.fields.marketboardProvider,
				itemProvider:        tt.fields.itemProvider,
			}
			got, err := profitProv.GetResaleInfoForItem(tt.args.ctx, tt.args.obj, tt.args.dataCenter, tt.args.homeServer, tt.args.buyCrystals, tt.args.buyFromOtherServers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResaleInfoForItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResaleInfoForItem() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemProfitProvider_GetCrossDcResaleProfit(t *testing.T) {
	type fields struct {
		maxValue            int
		returnUnlisted      bool
		recipeProvider      *database.RecipeDatabaseProvider
		marketboardProvider *database.MarketboardDatabaseProvider
		itemProvider        *database.ItemDatabaseProvider
	}
	type args struct {
		ctx        context.Context
		obj        *schema.Item
		dataCenter string
		homeServer string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *schema.ResaleInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profitProv := providers.ItemProfitProvider{
				maxValue:            tt.fields.maxValue,
				returnUnlisted:      tt.fields.returnUnlisted,
				recipeProvider:      tt.fields.recipeProvider,
				marketboardProvider: tt.fields.marketboardProvider,
				itemProvider:        tt.fields.itemProvider,
			}
			got, err := profitProv.GetCrossDcResaleProfit(tt.args.ctx, tt.args.obj, tt.args.dataCenter, tt.args.homeServer)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossDcResaleProfit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrossDcResaleProfit() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemProfitProvider_GetRecipePurchaseInfo(t *testing.T) {
	type fields struct {
		maxValue            int
		returnUnlisted      bool
		recipeProvider      *database.RecipeDatabaseProvider
		marketboardProvider *database.MarketboardDatabaseProvider
		itemProvider        *database.ItemDatabaseProvider
	}
	type args struct {
		componentItem       *schema.Item
		mbEntry             *schema.MarketboardEntry
		homeServer          string
		buyFromOtherServers *bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *schema.RecipePurchaseInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profitProv := providers.ItemProfitProvider{
				maxValue:            tt.fields.maxValue,
				returnUnlisted:      tt.fields.returnUnlisted,
				recipeProvider:      tt.fields.recipeProvider,
				marketboardProvider: tt.fields.marketboardProvider,
				itemProvider:        tt.fields.itemProvider,
			}
			if got := profitProv.GetRecipePurchaseInfo(tt.args.componentItem, tt.args.mbEntry, tt.args.homeServer, tt.args.buyFromOtherServers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRecipePurchaseInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemProfitProvider_GetVendorFlipProfit(t *testing.T) {
	type fields struct {
		maxValue            int
		returnUnlisted      bool
		recipeProvider      *database.RecipeDatabaseProvider
		marketboardProvider *database.MarketboardDatabaseProvider
		itemProvider        *database.ItemDatabaseProvider
	}
	type args struct {
		ctx        context.Context
		obj        *schema.Item
		dataCenter string
		homeServer string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profitProv := providers.ItemProfitProvider{
				maxValue:            tt.fields.maxValue,
				returnUnlisted:      tt.fields.returnUnlisted,
				recipeProvider:      tt.fields.recipeProvider,
				marketboardProvider: tt.fields.marketboardProvider,
				itemProvider:        tt.fields.itemProvider,
			}
			got, err := profitProv.GetVendorFlipProfit(tt.args.ctx, tt.args.obj, tt.args.dataCenter, tt.args.homeServer)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVendorFlipProfit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetVendorFlipProfit() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemProfitProvider_getHomeAndAwayItems(t *testing.T) {
	type fields struct {
		maxValue            int
		returnUnlisted      bool
		recipeProvider      *database.RecipeDatabaseProvider
		marketboardProvider *database.MarketboardDatabaseProvider
		itemProvider        *database.ItemDatabaseProvider
	}
	type args struct {
		marketEntry *schema.MarketboardEntry
		homeServer  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *schema.MarketEntry
		want1  *schema.MarketEntry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profitProv := providers.ItemProfitProvider{
				maxValue:            tt.fields.maxValue,
				returnUnlisted:      tt.fields.returnUnlisted,
				recipeProvider:      tt.fields.recipeProvider,
				marketboardProvider: tt.fields.marketboardProvider,
				itemProvider:        tt.fields.itemProvider,
			}
			got, got1 := profitProv.getHomeAndAwayItems(tt.args.marketEntry, tt.args.homeServer)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getHomeAndAwayItems() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getHomeAndAwayItems() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestItemProfitProvider_getItemValue(t *testing.T) {
	type fields struct {
		maxValue            int
		returnUnlisted      bool
		recipeProvider      *database.RecipeDatabaseProvider
		marketboardProvider *database.MarketboardDatabaseProvider
		itemProvider        *database.ItemDatabaseProvider
	}
	type args struct {
		marketEntries *schema.MarketboardEntry
		homeServer    string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profitProv := providers.ItemProfitProvider{
				maxValue:            tt.fields.maxValue,
				returnUnlisted:      tt.fields.returnUnlisted,
				recipeProvider:      tt.fields.recipeProvider,
				marketboardProvider: tt.fields.marketboardProvider,
				itemProvider:        tt.fields.itemProvider,
			}
			if got := profitProv.getItemValue(tt.args.marketEntries, tt.args.homeServer); got != tt.want {
				t.Errorf("getItemValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemProfitProvider_getRecipeResaleInfo(t *testing.T) {
	type fields struct {
		maxValue            int
		returnUnlisted      bool
		recipeProvider      *database.RecipeDatabaseProvider
		marketboardProvider *database.MarketboardDatabaseProvider
		itemProvider        *database.ItemDatabaseProvider
	}
	type args struct {
		ctx                 context.Context
		recipe              *schema.Recipe
		buyCrystals         *bool
		buyFromOtherServers *bool
		homeServer          string
		dataCenter          string
		itemValue           int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *schema.ResaleInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profitProv := providers.ItemProfitProvider{
				maxValue:            tt.fields.maxValue,
				returnUnlisted:      tt.fields.returnUnlisted,
				recipeProvider:      tt.fields.recipeProvider,
				marketboardProvider: tt.fields.marketboardProvider,
				itemProvider:        tt.fields.itemProvider,
			}
			got, err := profitProv.getRecipeResaleInfo(tt.args.ctx, tt.args.recipe, tt.args.buyCrystals, tt.args.buyFromOtherServers, tt.args.homeServer, tt.args.dataCenter, tt.args.itemValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRecipeResaleInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRecipeResaleInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewItemProfitProvider(t *testing.T) {
	type args struct {
		recipeProv *database.RecipeDatabaseProvider
		mbProv     *database.MarketboardDatabaseProvider
		itemProv   *database.ItemDatabaseProvider
	}
	tests := []struct {
		name string
		args args
		want *providers.ItemProfitProvider
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := providers.NewItemProfitProvider(tt.args.recipeProv, tt.args.mbProv, tt.args.itemProv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewItemProfitProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}
