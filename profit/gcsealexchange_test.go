package profitCalc

import (
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
	"testing"
)

func Test_CalculateSealValue(t *testing.T) {
	tests := []struct {
		name string
		item *readertype.Item
		want int
	}{
		{
			name: "iLvl 635 Returns correct value",
			item: &readertype.Item{
				Id:                 40766,
				Name:               "Voidmoon Coat of Fending",
				Description:        "",
				IconId:             57024,
				ItemLevel:          635,
				Rarity:             2,
				UiCategory:         35,
				SearchCategory:     0,
				SortCategory:       5,
				StackSize:          1,
				BuyFromVendorPrice: 0,
				SellToVendorPrice:  1219,
				ClassJobCategory:   59,
				CanBeTraded:        false,
				DropsFromDungeon:   true,
				CanBeHq:            false,
				IsCollectable:      false,
				IsGlamour:          true,
			},
			want: 1954,
		},
		{
			name: "iLvl 470 Returns correct value",
			item: &readertype.Item{
				Id:                 25608,
				Name:               "Augmented Deepshadow Earring of Fending",
				Description:        "",
				IconId:             55450,
				ItemLevel:          470,
				Rarity:             3,
				UiCategory:         41,
				SearchCategory:     0,
				SortCategory:       5,
				StackSize:          1,
				BuyFromVendorPrice: 0,
				SellToVendorPrice:  739,
				ClassJobCategory:   59,
				CanBeTraded:        false,
				DropsFromDungeon:   false,
				CanBeHq:            false,
				IsCollectable:      false,
				IsGlamour:          true,
			},
			want: 1673,
		},
		{
			name: "iLvl 375 Returns correct value",
			item: &readertype.Item{
				Id:                 24398,
				Name:               "Alliance Circlet of Healing",
				Description:        "",
				IconId:             41746,
				ItemLevel:          375,
				Rarity:             2,
				UiCategory:         34,
				SearchCategory:     0,
				SortCategory:       5,
				StackSize:          1,
				BuyFromVendorPrice: 0,
				SellToVendorPrice:  640,
				ClassJobCategory:   64,
				CanBeTraded:        false,
				DropsFromDungeon:   true,
				CanBeHq:            false,
				IsCollectable:      false,
				IsGlamour:          true,
			},
			want: 1500,
		},
		{
			name: "iLvl 205 Returns correct value",
			item: &readertype.Item{
				Id:                 14771,
				Name:               "Halone's Breeches of Fending",
				Description:        "",
				IconId:             45737,
				ItemLevel:          205,
				Rarity:             2,
				UiCategory:         36,
				SearchCategory:     0,
				SortCategory:       5,
				StackSize:          1,
				BuyFromVendorPrice: 0,
				SellToVendorPrice:  1005,
				ClassJobCategory:   59,
				CanBeTraded:        false,
				DropsFromDungeon:   false,
				CanBeHq:            false,
				IsCollectable:      false,
				IsGlamour:          true,
			},
			want: 1160,
		},
		{
			name: "iLvl 49 Returns correct value",
			item: &readertype.Item{
				Id:                 31571,
				Name:               "Aurum Boots",
				Description:        "",
				IconId:             49272,
				ItemLevel:          49,
				Rarity:             2,
				UiCategory:         38,
				SearchCategory:     0,
				SortCategory:       5,
				StackSize:          1,
				BuyFromVendorPrice: 0,
				SellToVendorPrice:  107,
				ClassJobCategory:   31,
				CanBeTraded:        false,
				DropsFromDungeon:   true,
				CanBeHq:            false,
				IsCollectable:      false,
				IsGlamour:          true,
			},
			want: 282,
		},
		{
			name: "Item without an iLvl Returns 0",
			item: &readertype.Item{
				Id:        1,
				Name:      "Non-existant Item",
				ItemLevel: 0,
				Rarity:    1,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := calculateSealValue(tt.item); got != tt.want {
					t.Errorf("calculateSealValue() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
