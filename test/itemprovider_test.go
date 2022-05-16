/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (itemprovider_test.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package test

import (
	schema "MarketMoogleAPI/core/graph/model"
	"MarketMoogleAPI/infrastructure/providers/db"
	"testing"
)

func TestItemProvider_GetItemFromApi(t *testing.T) {
	type fields struct {
		Db *db.DbProvider
	}

	type args struct {
		itemID int
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    schema.Item
		wantErr bool
	}{
		{
			name: "GetsFireCrystalCorrectly",
			fields: fields{
				Db: db.NewDbProvider(),
			},
			args: args{
				itemID: 8,
			},
			want: schema.Item{
				ItemID: 8,
				Name:   "Fire Crystal",
			},
			wantErr: false,
		},
		{
			name: "FailsAsExpected",
			fields: fields{
				Db: db.NewDbProvider(),
			},
			args: args{
				itemID: -1,
			},
			want:    schema.Item{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := db.ItemProvider{
				Db: tt.fields.Db,
			}
			got, err := i.GetItemFromApi(&tt.args.itemID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetItemFromApi() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && (got.ItemID != tt.want.ItemID && got.Name != tt.want.Name) {
				t.Errorf("GetItemFromApi() got = %v, want %v", got, tt.want)
			}
		})
	}
}
