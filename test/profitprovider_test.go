/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (profitprovider_test.go) is part of MarketMoogleAPI and is released under the GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package test

import (
	interfaces "MarketMoogleAPI/business/database"
	schema "MarketMoogleAPI/core/graph/model"
	"MarketMoogleAPI/core/util"
	"MarketMoogleAPI/infrastructure/providers"
	"MarketMoogleAPI/test/mocks"
	"context"
	"math"
	"reflect"
	"testing"
)

type TestProviders struct {
	recipeProvider      interfaces.RecipeProvider
	marketboardProvider interfaces.MarketBoardProvider
	itemProvider        interfaces.ItemProvider
}

func setupTestProviders() TestProviders {
	recipeProv := mocks.TestRecipeProvider{
		RecipeDatabase: make(map[int]*schema.Recipe),
	}

	mbProv := mocks.TestMarketboardProvider{
		MbDatabase: make(map[int]*schema.MarketboardEntry),
	}

	itemProv := mocks.TestItemProvider{
		ItemDatabase: make(map[int]*schema.Item),
	}

	return TestProviders{
		recipeProvider:      recipeProv,
		marketboardProvider: mbProv,
		itemProvider:        itemProv,
	}
}

func TestItemProfitProvider_GetCheapestOnDc(t *testing.T) {
	tests := []struct {
		name    string
		entries *schema.MarketboardEntry
		want    int
	}{
		{
			name: "No entries in DC returns max value",
			entries: &schema.MarketboardEntry{
				ItemID:              1,
				LastUpdateTime:      util.GetCurrentTimestampString(),
				MarketEntries:       []*schema.MarketEntry{},
				MarketHistory:       []*schema.MarketHistory{},
				DataCenter:          "Test Realm",
				CurrentAveragePrice: 0,
				CurrentMinPrice:     nil,
				RegularSaleVelocity: 0,
				HqSaleVelocity:      0,
				NqSaleVelocity:      0,
			},
			want: math.MaxInt32,
		},
		{
			name: "Returns cheapest correctly",
			entries: &schema.MarketboardEntry{
				ItemID:         1,
				LastUpdateTime: util.GetCurrentTimestampString(),
				MarketEntries: []*schema.MarketEntry{
					{
						ServerID:     1,
						Server:       "Cheap Server",
						Quantity:     1,
						TotalCost:    500,
						PricePer:     500,
						Hq:           false,
						IsCrafted:    false,
						RetainerName: util.MakePointer[string]("Cheap Seller"),
					},
					{
						ServerID:     2,
						Server:       "Expensive Server",
						Quantity:     1,
						TotalCost:    1000,
						PricePer:     1000,
						Hq:           false,
						IsCrafted:    false,
						RetainerName: util.MakePointer[string]("Expensive Seller"),
					},
				},
				MarketHistory:       []*schema.MarketHistory{},
				DataCenter:          "Test Realm",
				CurrentAveragePrice: 0,
				CurrentMinPrice:     nil,
				RegularSaleVelocity: 0,
				HqSaleVelocity:      0,
				NqSaleVelocity:      0,
			},
			want: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Other providers aren't necessary for this method
			profitProv := providers.NewItemProfitProvider(nil, nil, nil)

			if got := profitProv.GetCheapestOnDc(tt.entries); !reflect.DeepEqual(got.PricePer, tt.want) {
				t.Errorf("GetCheapestOnDc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemProfitProvider_GetCheapestOnServer(t *testing.T) {
	type args struct {
		entry  *schema.MarketboardEntry
		server string
	}
	tests := []struct {
		name     string
		args     args
		pricePer int
	}{
		{
			name: "Empty if no matching server found",
			args: args{
				entry: &schema.MarketboardEntry{
					ItemID:              1,
					LastUpdateTime:      util.GetCurrentTimestampString(),
					MarketEntries:       []*schema.MarketEntry{},
					MarketHistory:       []*schema.MarketHistory{},
					DataCenter:          "Test Realm",
					CurrentAveragePrice: 0,
					CurrentMinPrice:     nil,
					RegularSaleVelocity: 0,
					HqSaleVelocity:      0,
					NqSaleVelocity:      0,
				},
				server: "Test Server 1",
			},
			pricePer: math.MaxInt32,
		},
		{
			name: "Returns cheapest correctly",
			args: args{
				entry: &schema.MarketboardEntry{
					ItemID:         1,
					LastUpdateTime: util.GetCurrentTimestampString(),
					MarketEntries: []*schema.MarketEntry{
						{
							ServerID:     1,
							Server:       "Test Server 1",
							Quantity:     1,
							TotalCost:    500,
							PricePer:     500,
							Hq:           false,
							IsCrafted:    false,
							RetainerName: util.MakePointer[string]("Retainer 1"),
						},
						{
							ServerID:     2,
							Server:       "Test Server 1",
							Quantity:     1,
							TotalCost:    1000,
							PricePer:     1000,
							Hq:           false,
							IsCrafted:    false,
							RetainerName: util.MakePointer[string]("Retainer 2"),
						},
						{
							ServerID:     2,
							Server:       "Test Server 3",
							Quantity:     5,
							TotalCost:    500,
							PricePer:     100,
							Hq:           false,
							IsCrafted:    false,
							RetainerName: util.MakePointer[string]("Sellingway"),
						},
					},
					MarketHistory:       []*schema.MarketHistory{},
					DataCenter:          "Test Realm",
					CurrentAveragePrice: 0,
					CurrentMinPrice:     nil,
					RegularSaleVelocity: 0,
					HqSaleVelocity:      0,
					NqSaleVelocity:      0,
				},
				server: "Test Server 1",
			},
			pricePer: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Other providers aren't necessary for this method
			profitProv := providers.NewItemProfitProvider(nil, nil, nil)

			if got := profitProv.GetCheapestOnServer(tt.args.entry, tt.args.server); !reflect.DeepEqual(got.PricePer, tt.pricePer) {
				t.Errorf("GetCheapestOnServer() = %v, want %v", got, tt.pricePer)
			}
		})
	}
}

func TestItemProfitProvider_GetCrossDcResaleProfit(t *testing.T) {
	type args struct {
		obj        *schema.Item
		homeServer string
	}
	tests := []struct {
		name          string
		args          args
		marketEntries []*schema.MarketboardEntry
		want          *schema.ResaleInfo
		wantErr       bool
	}{
		{
			name: "Returns nil when no item presented",
			args: args{
				obj:        nil,
				homeServer: "",
			},
			marketEntries: nil,
			want:          nil,
			wantErr:       false,
		},
		{
			name: "Errors when item is not found",
			args: args{
				obj: &schema.Item{
					ItemID: 999,
					Name:   "Non-existant item",
					IconID: 999,
				},
				homeServer: "Non-existant land",
			},
			marketEntries: nil,
			want:          nil,
			wantErr:       true,
		},
		{
			name: "Returns error when item is not marketable",
			args: args{
				obj: &schema.Item{
					ItemID: 1,
					Name:   "Item with no market entries",
				},
				homeServer: "Marketless land",
			},
			marketEntries: nil,
			want: nil,
			wantErr: true,
		},
		{
			name: "Returns 0 profit when no market entries exist",
			args: args{
				obj: &schema.Item{
					ItemID: 1,
					Name:   "Item with no market entries",
				},
				homeServer: "Marketless land",
			},
			marketEntries: []*schema.MarketboardEntry{
				{
					ItemID:         1,
					LastUpdateTime: util.GetCurrentTimestampString(),
					MarketEntries: []*schema.MarketEntry{},
					MarketHistory:       nil,
					DataCenter:          "Test Data Center",
					CurrentAveragePrice: 0,
					CurrentMinPrice:     nil,
					RegularSaleVelocity: 0,
					HqSaleVelocity:      0,
					NqSaleVelocity:      0,
				},
			},
			want: &schema.ResaleInfo{
				Profit:          0,
				ItemID:          1,
				Quantity:        0,
				SingleCost:      math.MaxInt32,
				TotalCost:       0,
				ItemsToPurchase: []*schema.RecipePurchaseInfo{},
			},
			wantErr: false,
		},
		{
			name: "Correctly calculates cross-server flips",
			args: args{
				obj: &schema.Item{
					ItemID: 1,
					Name:   "Flippable item",
				},
				homeServer: "Home Server",
			},
			marketEntries: []*schema.MarketboardEntry{
				{
					ItemID:         1,
					LastUpdateTime: util.GetCurrentTimestampString(),
					MarketEntries: []*schema.MarketEntry{
						{
							ServerID:  1,
							Server:    "Home Server",
							Quantity:  1,
							TotalCost: 5000,
							PricePer:  5000,
						},
						{
							ServerID:  1,
							Server:    "Foreign Server",
							Quantity:  2,
							TotalCost: 2000,
							PricePer:  1000,
						},
						{
							ServerID:  2,
							Server:    "Foreign Server 2",
							Quantity:  2,
							TotalCost: 4000,
							PricePer:  2000,
						},
					},
					MarketHistory:       nil,
					DataCenter:          "Test Data Center",
					CurrentAveragePrice: 0,
					CurrentMinPrice:     nil,
					RegularSaleVelocity: 0,
					HqSaleVelocity:      0,
					NqSaleVelocity:      0,
				},
			},
			want: &schema.ResaleInfo{
				Profit:     8000,
				ItemID:     1,
				Quantity:   2,
				SingleCost: 1000,
				TotalCost:  2000,
				ItemsToPurchase: []*schema.RecipePurchaseInfo{
					{
						Item: &schema.Item{
							ItemID: 1,
							Name:   "Flippable item",
						},
						ServerToBuyFrom: "Foreign Server",
						BuyFromVendor:   false,
						PricePer:        1000,
						TotalCost:       2000,
						Quantity:        2,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Correctly calculates same-server flips",
			args: args{
				obj: &schema.Item{
					ItemID: 1,
					Name:   "Flippable item",
				},
				homeServer: "Home Server",
			},
			marketEntries: []*schema.MarketboardEntry{
				{
					ItemID:         1,
					LastUpdateTime: util.GetCurrentTimestampString(),
					MarketEntries: []*schema.MarketEntry{
						{
							ServerID:  1,
							Server:    "Home Server",
							Quantity:  1,
							TotalCost: 500,
							PricePer:  500,
						},
						{
							ServerID:  1,
							Server:    "Home Server",
							Quantity:  1,
							TotalCost: 5000,
							PricePer:  5000,
						},
						{
							ServerID:  1,
							Server:    "Home Server",
							Quantity:  1,
							TotalCost: 6000,
							PricePer:  6000,
						},
						{
							ServerID:  2,
							Server:    "Foreign Server",
							Quantity:  2,
							TotalCost: 980,
							PricePer:  490,
						},
					},
					MarketHistory:       nil,
					DataCenter:          "Test Data Center",
					CurrentAveragePrice: 0,
					CurrentMinPrice:     nil,
					RegularSaleVelocity: 0,
					HqSaleVelocity:      0,
					NqSaleVelocity:      0,
				},
			},
			want: &schema.ResaleInfo{
				Profit:     4500,
				ItemID:     1,
				Quantity:   1,
				SingleCost: 500,
				TotalCost:  500,
				ItemsToPurchase: []*schema.RecipePurchaseInfo{
					{
						Item: &schema.Item{
							ItemID: 1,
							Name:   "Flippable item",
						},
						ServerToBuyFrom: "Home Server",
						BuyFromVendor:   false,
						PricePer:        500,
						TotalCost:       500,
						Quantity:        1,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			testProviders := setupTestProviders()

			profitProv := providers.NewItemProfitProvider(
				testProviders.recipeProvider,
				testProviders.marketboardProvider,
				testProviders.itemProvider)

			//Add market entries to provider for test
			for _, marketEntry := range tt.marketEntries {
				_, err := testProviders.marketboardProvider.CreateMarketEntry(ctx, marketEntry)

				if err != nil {
					t.Errorf("error adding market entries to test provider")
				}
			}

			got, err := profitProv.GetCrossDcResaleProfit(ctx, tt.args.obj, "Test Data Center", tt.args.homeServer)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossDcResaleProfit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !resaleInfoObjectsAreSame(got, tt.want) {
				t.Errorf("GetCrossDcResaleProfit() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func resaleInfoObjectsAreSame(got *schema.ResaleInfo, want *schema.ResaleInfo) bool {
	if want == nil || got == nil {
		return true
	}

	//Since we're not creating the same pointer for object comparison we will have to compare elements individually
	sameProperties := true
	if got.ItemID != want.ItemID {
		sameProperties = false
	}

	if got.SingleCost != want.SingleCost {
		sameProperties = false
	}

	if got.TotalCost != want.TotalCost {
		sameProperties = false
	}

	if got.Quantity != want.Quantity {
		sameProperties = false
	}

	if got.Profit != want.Profit {
		sameProperties = false
	}

	for index := range want.ItemsToPurchase {
		if !reflect.DeepEqual(*got.ItemsToPurchase[index], *want.ItemsToPurchase[index]) {
			sameProperties = false
			break
		}
	}

	return sameProperties
}

/*
func TestItemProfitProvider_GetRecipePurchaseInfo(t *testing.T) {
	type fields struct {
		maxValue            int
		returnUnlisted      bool
		recipeProvider      database.RecipeProvider
		marketboardProvider database.MarketBoardProvider
		itemProvider        database.ItemProvider
	}
	type args struct {
		componentItem       *schema.Item
		mbEntry             *schema.MarketboardEntry
		homeServer          string
		buyFromOtherServers *bool
		count               int
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
			if got := profitProv.GetRecipePurchaseInfo(tt.args.componentItem, tt.args.mbEntry, tt.args.homeServer, tt.args.buyFromOtherServers, tt.args.count); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRecipePurchaseInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemProfitProvider_GetResaleInfoForItem(t *testing.T) {
	type fields struct {
		maxValue            int
		returnUnlisted      bool
		recipeProvider      database.RecipeProvider
		marketboardProvider database.MarketBoardProvider
		itemProvider        database.ItemProvider
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

func TestItemProfitProvider_GetVendorFlipProfit(t *testing.T) {
	type fields struct {
		maxValue            int
		returnUnlisted      bool
		recipeProvider      database.RecipeProvider
		marketboardProvider database.MarketBoardProvider
		itemProvider        database.ItemProvider
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
		recipeProvider      database.RecipeProvider
		marketboardProvider database.MarketBoardProvider
		itemProvider        database.ItemProvider
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
		recipeProvider      database.RecipeProvider
		marketboardProvider database.MarketBoardProvider
		itemProvider        database.ItemProvider
	}
	type args struct {
		marketEntries *schema.MarketboardEntry
		server        string
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
			if got := profitProv.getItemValue(tt.args.marketEntries, tt.args.server); got != tt.want {
				t.Errorf("getItemValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemProfitProvider_getRecipeResaleInfo(t *testing.T) {
	type fields struct {
		maxValue            int
		returnUnlisted      bool
		recipeProvider      database.RecipeProvider
		marketboardProvider database.MarketBoardProvider
		itemProvider        database.ItemProvider
	}
	type args struct {
		ctx                 context.Context
		recipe              *schema.Recipe
		buyCrystals         *bool
		buyFromOtherServers *bool
		homeServer          string
		dataCenter          string
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
			got, err := profitProv.getRecipeResaleInfo(tt.args.ctx, tt.args.recipe, tt.args.buyCrystals, tt.args.buyFromOtherServers, tt.args.homeServer, tt.args.dataCenter)
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
*/
