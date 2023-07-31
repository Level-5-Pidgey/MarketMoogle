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

const TestDataCenter string = "Test Data Center"
const HomeServer string = "Home Server"
const AwayServer string = "Foreign Server"

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

func recipeResaleObjectsAreTheSame(got *schema.RecipeProfitInfo, want *schema.RecipeProfitInfo) bool {
	if want == nil || got == nil {
		return true
	}

	//Since we're not creating the same pointer for object comparison we will have to compare elements individually
	if got.CraftType != want.CraftType {
		return false
	}

	if got.CraftLevel != want.CraftLevel {
		return false
	}

	if (got.ResaleInfo == nil && want.ResaleInfo != nil) || (got.ResaleInfo != nil && want.ResaleInfo == nil) {
		return false
	}

	if got.ResaleInfo != nil && want.ResaleInfo != nil {
		return resaleInfoObjectsAreSame(got.ResaleInfo, want.ResaleInfo)
	}

	return true
}

func resaleInfoObjectsAreSame(got *schema.ProfitInfo, want *schema.ProfitInfo) bool {
	if want == nil || got == nil {
		return true
	}

	if got.ItemID != want.ItemID {
		return false
	}

	if got.SingleCost != want.SingleCost {
		return false
	}

	if got.TotalCost != want.TotalCost {
		return false
	}

	if got.Quantity != want.Quantity {
		return false
	}

	if got.Profit != want.Profit {
		return false
	}

	for index := range want.ItemsToPurchase {
		if !reflect.DeepEqual(*got.ItemsToPurchase[index], *want.ItemsToPurchase[index]) {
			return false
		}
	}

	return true
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

			if got := profitProv.GetCheapestOnDataCenter(tt.entries); !reflect.DeepEqual(got.TotalCost, tt.want) {
				t.Errorf("GetCheapestOnDataCenter() = %v, want %v", got, tt.want)
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
		want          *schema.ProfitInfo
		wantErr       bool
	}{
		{
			name: "Returns nil when no items presented",
			args: args{
				obj:        nil,
				homeServer: "",
			},
			marketEntries: nil,
			want:          nil,
			wantErr:       false,
		},
		{
			name: "Errors when items is not found",
			args: args{
				obj: &schema.Item{
					Id:     999,
					Name:   "Non-existant items",
					IconID: 999,
				},
				homeServer: "Non-existant land",
			},
			marketEntries: nil,
			want:          nil,
			wantErr:       true,
		},
		{
			name: "Returns error when items is not marketable",
			args: args{
				obj: &schema.Item{
					Id:   1,
					Name: "Item with no market entries",
				},
				homeServer: "Marketless land",
			},
			marketEntries: nil,
			want:          nil,
			wantErr:       true,
		},
		{
			name: "Returns 0 profit when no market entries exist",
			args: args{
				obj: &schema.Item{
					Id:   1,
					Name: "Item with no market entries",
				},
				homeServer: "Marketless land",
			},
			marketEntries: []*schema.MarketboardEntry{
				{
					ItemID:         1,
					LastUpdateTime: util.GetCurrentTimestampString(),
					MarketEntries:  []*schema.MarketEntry{},
					DataCenter:     TestDataCenter,
				},
			},
			want: &schema.ProfitInfo{
				Profit:          0,
				ItemID:          1,
				Quantity:        0,
				SingleCost:      math.MaxInt32,
				TotalCost:       math.MaxInt32,
				ItemsToPurchase: []*schema.ItemCostInfo{},
			},
			wantErr: false,
		},
		{
			name: "Correctly calculates cross-server flips",
			args: args{
				obj: &schema.Item{
					Id:   1,
					Name: "Flippable items",
				},
				homeServer: HomeServer,
			},
			marketEntries: []*schema.MarketboardEntry{
				{
					ItemID:         1,
					LastUpdateTime: util.GetCurrentTimestampString(),
					MarketEntries: []*schema.MarketEntry{
						{
							ServerID:  1,
							Server:    HomeServer,
							Quantity:  1,
							TotalCost: 5000,
							PricePer:  5000,
						},
						{
							ServerID:  1,
							Server:    AwayServer,
							Quantity:  2,
							TotalCost: 2000,
							PricePer:  1000,
						},
						{
							ServerID:  2,
							Server:    AwayServer + "2",
							Quantity:  2,
							TotalCost: 4000,
							PricePer:  2000,
						},
					},
					DataCenter: TestDataCenter,
				},
			},
			want: &schema.ProfitInfo{
				Profit:     8000,
				ItemID:     1,
				Quantity:   2,
				SingleCost: 1000,
				TotalCost:  2000,
				ItemsToPurchase: []*schema.ItemCostInfo{
					{
						Item: &schema.Item{
							Id:   1,
							Name: "Flippable items",
						},
						ServerToBuyFrom: AwayServer,
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
					Id:   1,
					Name: "Flippable items",
				},
				homeServer: HomeServer,
			},
			marketEntries: []*schema.MarketboardEntry{
				{
					ItemID:         1,
					LastUpdateTime: util.GetCurrentTimestampString(),
					MarketEntries: []*schema.MarketEntry{
						{
							ServerID:  1,
							Server:    HomeServer,
							Quantity:  1,
							TotalCost: 500,
							PricePer:  500,
						},
						{
							ServerID:  1,
							Server:    HomeServer,
							Quantity:  1,
							TotalCost: 5000,
							PricePer:  5000,
						},
						{
							ServerID:  1,
							Server:    HomeServer,
							Quantity:  1,
							TotalCost: 6000,
							PricePer:  6000,
						},
						{
							ServerID:  2,
							Server:    AwayServer,
							Quantity:  2,
							TotalCost: 980,
							PricePer:  490,
						},
					},
					DataCenter: TestDataCenter,
				},
			},
			want: &schema.ProfitInfo{
				Profit:     4500,
				ItemID:     1,
				Quantity:   1,
				SingleCost: 500,
				TotalCost:  500,
				ItemsToPurchase: []*schema.ItemCostInfo{
					{
						Item: &schema.Item{
							Id:   1,
							Name: "Flippable items",
						},
						ServerToBuyFrom: HomeServer,
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

			got, err := profitProv.GetCrossDcResaleProfit(ctx, tt.args.obj, TestDataCenter, tt.args.homeServer)

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

func TestItemProfitProvider_GetRecipePurchaseInfo(t *testing.T) {
	type args struct {
		item                *schema.Item
		mbEntry             *schema.MarketboardEntry
		homeServer          string
		buyFromOtherServers *bool
		count               int
	}
	tests := []struct {
		name string
		args args
		want *schema.ItemCostInfo
	}{
		{
			name: "Always buys from vendor if possible",
			args: args{
				item: &schema.Item{
					Id:                 1,
					Name:               "Test Item",
					SellToVendorValue:  util.MakePointer[int](50),
					BuyFromVendorValue: util.MakePointer[int](500),
				},
				mbEntry: &schema.MarketboardEntry{
					ItemID:         1,
					LastUpdateTime: util.GetCurrentTimestampString(),
					MarketEntries: []*schema.MarketEntry{
						{
							ServerID:  1,
							Server:    HomeServer,
							Quantity:  1,
							TotalCost: 400,
							PricePer:  400,
						},
					},
					DataCenter: TestDataCenter,
				},
				homeServer:          HomeServer,
				buyFromOtherServers: nil,
				count:               1,
			},
			want: &schema.ItemCostInfo{
				Item: &schema.Item{
					Id:                 1,
					Name:               "Test Item",
					SellToVendorValue:  util.MakePointer[int](50),
					BuyFromVendorValue: util.MakePointer[int](500),
				},
				ServerToBuyFrom: HomeServer,
				BuyFromVendor:   true,
				PricePer:        500,
				TotalCost:       500,
				Quantity:        1,
			},
		},
		{
			name: "Buys from cheapest server on the DC",
			args: args{
				item: &schema.Item{
					Id:   1,
					Name: "Test Item",
				},
				mbEntry: &schema.MarketboardEntry{
					ItemID:         1,
					LastUpdateTime: util.GetCurrentTimestampString(),
					MarketEntries: []*schema.MarketEntry{
						{
							ServerID:  2,
							Server:    AwayServer,
							Quantity:  1,
							TotalCost: 400,
							PricePer:  400,
						},
						{
							ServerID:  2,
							Server:    AwayServer,
							Quantity:  2,
							TotalCost: 600,
							PricePer:  300,
						},
						{
							ServerID:  2,
							Server:    AwayServer,
							Quantity:  1,
							TotalCost: 450,
							PricePer:  450,
						},
						{
							ServerID:  1,
							Server:    HomeServer,
							Quantity:  2,
							TotalCost: 10000,
							PricePer:  5000,
						},
					},
					DataCenter: TestDataCenter,
				},
				homeServer:          HomeServer,
				buyFromOtherServers: util.MakePointer[bool](true),
				count:               1,
			},
			want: &schema.ItemCostInfo{
				Item: &schema.Item{
					Id:   1,
					Name: "Test Item",
				},
				ServerToBuyFrom: AwayServer,
				BuyFromVendor:   false,
				PricePer:        400,
				TotalCost:       400,
				Quantity:        1,
			},
		},
		{
			name: "Buys from cheapest at home if option toggled",
			args: args{
				item: &schema.Item{
					Id:   1,
					Name: "Test Item",
				},
				mbEntry: &schema.MarketboardEntry{
					ItemID:         1,
					LastUpdateTime: util.GetCurrentTimestampString(),
					MarketEntries: []*schema.MarketEntry{
						{
							ServerID:  1,
							Server:    HomeServer,
							Quantity:  1,
							TotalCost: 400,
							PricePer:  400,
						},
						{
							ServerID:  1,
							Server:    HomeServer,
							Quantity:  2,
							TotalCost: 700,
							PricePer:  350,
						},
						{
							ServerID:  2,
							Server:    AwayServer,
							Quantity:  3,
							TotalCost: 900,
							PricePer:  300,
						},
						{
							ServerID:  2,
							Server:    AwayServer,
							Quantity:  1,
							TotalCost: 300,
							PricePer:  300,
						},
					},
					DataCenter: TestDataCenter,
				},
				homeServer:          HomeServer,
				buyFromOtherServers: util.MakePointer[bool](false),
				count:               1,
			},
			want: &schema.ItemCostInfo{
				Item: &schema.Item{
					Id:   1,
					Name: "Test Item",
				},
				ServerToBuyFrom: HomeServer,
				BuyFromVendor:   false,
				PricePer:        400,
				TotalCost:       400,
				Quantity:        1,
			},
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
			_, err := testProviders.marketboardProvider.CreateMarketEntry(ctx, tt.args.mbEntry)
			if err != nil {
				t.Errorf("error adding market entries to test provider")
			}

			//Add basic items to items provider
			_, err = testProviders.itemProvider.InsertItem(ctx, tt.args.item)
			if err != nil {
				t.Errorf("error adding test items to provider")
			}

			if got := profitProv.GetComponentCostInfo(tt.args.item, tt.args.mbEntry, tt.args.homeServer, tt.args.buyFromOtherServers, tt.args.count); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetComponentCostInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemProfitProvider_GetRecipeProfitForItem(t *testing.T) {
	type args struct {
		items               []*schema.Item
		marketboardEntries  []*schema.MarketboardEntry
		recipe              *schema.Recipe
		dataCenter          string
		homeServer          string
		buyCrystals         *bool
		buyFromOtherServers *bool
	}
	tests := []struct {
		name    string
		args    args
		want    *schema.RecipeProfitInfo
		wantErr bool
	}{
		{
			name: "Returns early if items doesn't exist",
			args: args{
				items: []*schema.Item{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Returns early if items has no recipe",
			args: args{
				items: []*schema.Item{
					{
						Id:   1,
						Name: "Recipe-less items",
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Correctly identifies cost of recipe",
			args: args{
				items: []*schema.Item{
					{
						Id:   20,
						Name: "Tasty Cake",
					},
					{
						Id:   21,
						Name: "Cake Mix",
					},
					{
						Id:   22,
						Name: "Eggs",
					},
				},
				marketboardEntries: []*schema.MarketboardEntry{
					{
						ItemID:         20,
						LastUpdateTime: util.GetCurrentTimestampString(),
						MarketEntries: []*schema.MarketEntry{
							{
								ServerID:  1,
								Server:    HomeServer,
								Quantity:  2,
								TotalCost: 1000,
								PricePer:  500,
							},
						},
						DataCenter:          TestDataCenter,
						CurrentAveragePrice: 500,
					},
					{
						ItemID:         21,
						LastUpdateTime: util.GetCurrentTimestampString(),
						MarketEntries: []*schema.MarketEntry{
							{
								ServerID:  2,
								Server:    AwayServer,
								Quantity:  1,
								TotalCost: 250,
								PricePer:  250,
							},
						},
						DataCenter:          TestDataCenter,
						CurrentAveragePrice: 250,
					},
					{
						ItemID:         22,
						LastUpdateTime: util.GetCurrentTimestampString(),
						MarketEntries: []*schema.MarketEntry{
							{
								ServerID:  1,
								Server:    HomeServer,
								Quantity:  4,
								TotalCost: 400,
								PricePer:  100,
							},
						},
						DataCenter:          TestDataCenter,
						CurrentAveragePrice: 100,
					},
				},
				recipe: &schema.Recipe{
					RecipeID:       1,
					ItemResultID:   20,
					ResultQuantity: 1,
					CraftedBy:      schema.CrafterTypeCulinarian,
					RecipeLevel:    util.MakePointer[int](50),
					RecipeItems: []*schema.RecipeContents{
						{
							ItemID: 21,
							Count:  1,
						},
						{
							ItemID: 22,
							Count:  2,
						},
					},
				},
				dataCenter:          TestDataCenter,
				homeServer:          HomeServer,
				buyCrystals:         nil,
				buyFromOtherServers: nil,
			},
			want: &schema.RecipeProfitInfo{
				ResaleInfo: &schema.ProfitInfo{
					Profit:     50,
					ItemID:     20,
					Quantity:   1,
					SingleCost: 450,
					TotalCost:  450,
					ItemsToPurchase: []*schema.ItemCostInfo{
						{
							Item: &schema.Item{
								Id:   21,
								Name: "Cake Mix",
							},
							ServerToBuyFrom: AwayServer,
							BuyFromVendor:   false,
							PricePer:        250,
							TotalCost:       250,
							Quantity:        1,
						},
						{
							Item: &schema.Item{
								Id:   22,
								Name: "Eggs",
							},
							ServerToBuyFrom: HomeServer,
							BuyFromVendor:   false,
							PricePer:        100,
							TotalCost:       400,
							Quantity:        4,
						},
					},
				},
				CraftLevel: 50,
				CraftType:  schema.CrafterTypeCulinarian,
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
			for _, mbEntry := range tt.args.marketboardEntries {
				_, err := testProviders.marketboardProvider.CreateMarketEntry(ctx, mbEntry)

				if err != nil {
					t.Errorf("error adding market entries to test provider")
				}
			}

			//Add items to provider
			for _, item := range tt.args.items {
				_, err := testProviders.itemProvider.InsertItem(ctx, item)

				if err != nil {
					t.Errorf("error adding test items to provider")
				}
			}

			//Add items recipe to recipe provider
			_, err := testProviders.recipeProvider.InsertRecipe(ctx, tt.args.recipe)
			if err != nil {
				t.Errorf("error adding test recipe to provider")
			}

			//Return early if there's no items
			if len(tt.args.items) == 0 {
				return
			}

			got, err := profitProv.GetRecipeProfitForItem(ctx, tt.args.items[0], tt.args.dataCenter, tt.args.homeServer, tt.args.buyCrystals, tt.args.buyFromOtherServers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecipeProfitForItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !recipeResaleObjectsAreTheSame(got, tt.want) {
				t.Errorf("GetRecipeProfitForItem() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemProfitProvider_GetVendorFlipProfit(t *testing.T) {
	type args struct {
		marketboardEntries []*schema.MarketboardEntry
		obj                *schema.Item
		dataCenter         string
		homeServer         string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Returns 0 profit if there's no object provided",
			args: args{
				marketboardEntries: []*schema.MarketboardEntry{
					{
						ItemID:         1,
						LastUpdateTime: util.GetCurrentTimestampString(),
						MarketEntries: []*schema.MarketEntry{
							{
								ServerID:  1,
								Server:    HomeServer,
								Quantity:  1,
								TotalCost: 5000,
								PricePer:  5000,
							},
						},
						DataCenter: TestDataCenter,
					},
				},
				obj:        nil,
				dataCenter: TestDataCenter,
				homeServer: HomeServer,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "Errors if there's no marketboard entries found for the object",
			args: args{
				marketboardEntries: []*schema.MarketboardEntry{
					{
						ItemID:         2,
						LastUpdateTime: util.GetCurrentTimestampString(),
						MarketEntries: []*schema.MarketEntry{
							{
								ServerID:  1,
								Server:    HomeServer,
								Quantity:  1,
								TotalCost: 5000,
								PricePer:  5000,
							},
						},
						DataCenter: TestDataCenter,
					},
				},
				obj: &schema.Item{
					Id:                 2,
					Name:               "Unmarketable Item",
					BuyFromVendorValue: util.MakePointer(5000),
				},
				dataCenter: TestDataCenter,
				homeServer: HomeServer,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "Returns 0 if the object can't be bought from a vendor",
			args: args{
				marketboardEntries: []*schema.MarketboardEntry{
					{
						ItemID:         2,
						LastUpdateTime: util.GetCurrentTimestampString(),
						MarketEntries: []*schema.MarketEntry{
							{
								ServerID:  1,
								Server:    HomeServer,
								Quantity:  1,
								TotalCost: 5000,
								PricePer:  5000,
							},
						},
						DataCenter: TestDataCenter,
					},
				},
				obj: &schema.Item{
					Id:                 1,
					Name:               "Unmarketable Item",
					BuyFromVendorValue: nil,
				},
				dataCenter: TestDataCenter,
				homeServer: HomeServer,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "Correctly identifies profit for items that can be bought from a vendor and are sold on the home server",
			args: args{
				marketboardEntries: []*schema.MarketboardEntry{
					{
						ItemID:         1,
						LastUpdateTime: util.GetCurrentTimestampString(),
						MarketEntries: []*schema.MarketEntry{
							{
								ServerID:  1,
								Server:    HomeServer,
								Quantity:  1,
								TotalCost: 5000,
								PricePer:  5000,
							},
						},
						DataCenter: TestDataCenter,
					},
				},
				obj: &schema.Item{
					Id:                 1,
					Name:               "Marketable Item",
					BuyFromVendorValue: util.MakePointer(3000),
				},
				dataCenter: TestDataCenter,
				homeServer: HomeServer,
			},
			want:    2000,
			wantErr: false,
		},
		{
			name: "Reports 0 profit if there are no home server market entries",
			args: args{
				marketboardEntries: []*schema.MarketboardEntry{
					{
						ItemID:         1,
						LastUpdateTime: util.GetCurrentTimestampString(),
						MarketEntries: []*schema.MarketEntry{
							{
								ServerID:  1,
								Server:    AwayServer,
								Quantity:  1,
								TotalCost: 5000,
								PricePer:  5000,
							},
						},
						DataCenter: TestDataCenter,
					},
				},
				obj: &schema.Item{
					Id:                 1,
					Name:               "Marketable Item",
					BuyFromVendorValue: util.MakePointer(3000),
				},
				dataCenter: TestDataCenter,
				homeServer: HomeServer,
			},
			want:    0,
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
			for _, mbEntry := range tt.args.marketboardEntries {
				_, err := testProviders.marketboardProvider.CreateMarketEntry(ctx, mbEntry)

				if err != nil {
					t.Errorf("error adding market entries to test provider")
				}
			}

			//Add obj to provider
			_, err := testProviders.itemProvider.InsertItem(ctx, tt.args.obj)

			if err != nil {
				t.Errorf("error adding test items to provider")
			}

			got, err := profitProv.GetVendorFlipProfit(ctx, tt.args.obj, tt.args.dataCenter, tt.args.homeServer)
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
