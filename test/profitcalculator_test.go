/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (profitcalculator_test.go) is part of MarketMoogleAPI and is released under the GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package test

import (
	schema "MarketMoogleAPI/core/graph/model"
	"MarketMoogleAPI/infrastructure/providers"
	"context"
	"reflect"
	"testing"
)

func TestItemProfitProvider_GetCheapestOnServer(t *testing.T) {
	type args struct {
		entry  schema.MarketboardEntry
		server string
	}
	var tests = []struct {
		name string
		args args
		want int
	}{
		{
			name: "",
			args: args{
				entry: schema.MarketboardEntry{
					ItemID:         14693,
					LastUpdateTime: "1657512123755",
					MarketEntries: []*schema.MarketEntry{
						{
							ServerID:     1,
							Server:       "TestServer",
							Quantity:     1,
							TotalPrice:   3500,
							PricePer:     3500,
							Hq:           false,
							IsCrafted:    true,
							RetainerName: nil,
						},
					},
					MarketHistory:       nil,
					DataCenter:          "TestDC",
					CurrentAveragePrice: 3500,
					CurrentMinPrice:     nil,
					RegularSaleVelocity: 0.1431,
					HqSaleVelocity:      0,
					NqSaleVelocity:      0.1431,
				},
				server: "TestServer",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profitProv := providers.NewItemProfitProvider()

			if got := profitProv.GetCheapestOnServer(&tt.args.entry, &tt.args.server); got != tt.want {
				t.Errorf("GetCheapestOnServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemProfitProvider_GetCheapestPriceAndServer(t *testing.T) {
	type args struct {
		entry *schema.MarketboardEntry
	}
	tests := []struct {
		name   string
		args   args
		want   int
		want1  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profitProv := providers.NewItemProfitProvider()
			
			got, got1 := profitProv.GetCheapestPriceAndServer(tt.args.entry)
			
			if got != tt.want {
				t.Errorf("GetCheapestPriceAndServer() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetCheapestPriceAndServer() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestItemProfitProvider_GetCraftingProfit(t *testing.T) {
	type args struct {
		ctx                 context.Context
		obj                 *schema.Item
		dataCenter          *string
		homeServer          *string
		buyCrystals         *bool
		buyFromOtherServers *bool
	}
	tests := []struct {
		name    string
		args    args
		want    *schema.RecipeResaleInformation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profitProv := providers.NewItemProfitProvider()
			
			got, err := profitProv.GetCraftingProfit(tt.args.ctx, tt.args.obj, tt.args.dataCenter, tt.args.homeServer, tt.args.buyCrystals, tt.args.buyFromOtherServers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCraftingProfit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCraftingProfit() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemProfitProvider_GetCrossDcResaleProfit(t *testing.T) {
	type args struct {
		ctx        context.Context
		obj        *schema.Item
		dataCenter *string
		homeServer *string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profitProv := providers.NewItemProfitProvider()
			
			got, err := profitProv.GetCrossDcResaleProfit(tt.args.ctx, tt.args.obj, tt.args.dataCenter, tt.args.homeServer)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossDcResaleProfit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCrossDcResaleProfit() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemProfitProvider_GetItemValue(t *testing.T) {
	type args struct {
		componentItem          *schema.Item
		itemMarketboardEntries *schema.MarketboardEntry
		homeServer             *string
		buyFromOtherServers    *bool
	}
	tests := []struct {
		name   string
		args   args
		want   *int
		want1  *string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profitProv := providers.NewItemProfitProvider()
			
			got, got1 := profitProv.GetItemValue(tt.args.componentItem, tt.args.itemMarketboardEntries, tt.args.homeServer, tt.args.buyFromOtherServers)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItemValue() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetItemValue() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestItemProfitProvider_GetVendorResaleProfit(t *testing.T) {
	type args struct {
		ctx        context.Context
		obj        *schema.Item
		dataCenter *string
		homeServer *string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profitProv := providers.NewItemProfitProvider()
			
			got, err := profitProv.GetVendorResaleProfit(tt.args.ctx, tt.args.obj, tt.args.dataCenter, tt.args.homeServer)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVendorResaleProfit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetVendorResaleProfit() got = %v, want %v", got, tt.want)
			}
		})
	}
}
