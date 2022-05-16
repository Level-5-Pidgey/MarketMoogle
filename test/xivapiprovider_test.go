/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (xivapiprovider_test.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package test

import (
	"MarketMoogleAPI/core/apitypes/xivapi"
	"MarketMoogleAPI/infrastructure/providers"
	"reflect"
	"testing"
)

func TestXivApiProvider_GetGameItemById(t *testing.T) {
	type args struct {
		contentId int
	}
	tests := []struct {
		name    string
		args    args
		want    xivapi.GameItem
		wantErr bool
	}{
		{
			name: "GetsEarthCrystalCorrectly",
			args: args{
				contentId: 11,
			},
			want: xivapi.GameItem{
				ID:   11,
				Name: "Earth Crystal",
			},
			wantErr: false,
		},
		{
			name: "HandlesApiFailure",
			args: args{
				contentId: -1,
			},
			want:    xivapi.GameItem{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := providers.XivApiProvider{}
			got, err := p.GetGameItemById(&tt.args.contentId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGameItemById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				gotDereferenced := *got
				if gotDereferenced.ID != tt.want.ID || gotDereferenced.Name != tt.want.Name {
					t.Errorf("GetGameItemById() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestXivApiProvider_GetItemsAndPrices(t *testing.T) {
	type args struct {
		shopId int
	}
	tests := []struct {
		name    string
		args    args
		want    map[int]int
		wantErr bool
	}{
		{
			name: "GetsMinionShopItems",
			args: args{
				shopId: 262574,
			},
			want:    map[int]int{6003: 2400, 6004: 2400, 6005: 2400},
			wantErr: false,
		},
		{
			name: "HandlesErrorCorrectly",
			args: args{
				shopId: -1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := providers.XivApiProvider{}
			got, err := p.GetItemsAndPrices(&tt.args.shopId)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetItemsAndPrices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("GetItemsAndPrices() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXivApiProvider_GetItemsInShop(t *testing.T) {
	type args struct {
		shopId int
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "GetsMinionGilShopItems",
			args: args{
				shopId: 262574,
			},
			want: []int{
				6003,
				6004,
				6005,
			},
			wantErr: false,
		},
		{
			name: "HandlesErrorCorrectly",
			args: args{
				shopId: -1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := providers.XivApiProvider{}
			got, err := p.GetItemsAndPrices(&tt.args.shopId)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetItemsAndPrices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("GetItemsAndPrices() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXivApiProvider_GetLodestoneInfoById(t *testing.T) {
	type args struct {
		lodestoneId int
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "GetPidgeyLodestoneName",
			args: args{
				lodestoneId: 41253877,
			},
			want:    "Pijii Otto",
			wantErr: false,
		},
		{
			name: "GetSarahLodestoneName",
			args: args{
				lodestoneId: 11927482,
			},
			want:    "S'arah Jane",
			wantErr: false,
		},
		{
			name: "GetEmelyLodestoneName",
			args: args{
				lodestoneId: 0,
			},
			want:    "Phrawg Xhula",
			wantErr: false,
		},
		{
			name: "GetNishimbaLodestoneName",
			args: args{
				lodestoneId: 24789452,
			},
			want:    "Nishimba Orion",
			wantErr: false,
		},
		{
			name: "GetBadLodestoneName",
			args: args{
				lodestoneId: -1,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := providers.XivApiProvider{}
			got, err := p.GetLodestoneInfoById(&tt.args.lodestoneId)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLodestoneInfoById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.want != got.Character.Name {
				t.Errorf("GetLodestoneInfoById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXivApiProvider_GetRecipeIdByItemId(t *testing.T) {
	type args struct {
		contentId int
	}

	type recipeLookupIds struct {
		armTargetId int
		bsmTargetId int
	}

	tests := []struct {
		name    string
		args    args
		want    recipeLookupIds
		wantErr bool
	}{
		{
			name: "GetMythrilRivetsRecipeCorrectly",
			args: args{
				contentId: 5099,
			},
			want: recipeLookupIds{
				armTargetId: 263,
				bsmTargetId: 139,
			},
			wantErr: false,
		},
		{
			name: "HandlesErrorCorrectly",
			args: args{
				contentId: -1,
			},
			want:    recipeLookupIds{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := providers.XivApiProvider{}
			got, err := p.GetRecipeIdByItemId(&tt.args.contentId)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecipeIdByItemId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && (tt.want.armTargetId != got.ARMTargetID || tt.want.bsmTargetId != got.BSMTargetID) {
				t.Errorf("GetRecipeIdByItemId() got = %v, want %v", got, tt.want)
			}
		})
	}
}
