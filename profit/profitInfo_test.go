package profitCalc

import (
	"github.com/level-5-pidgey/MarketMoogleApi/db"
	"reflect"
	"testing"
)

func TestProfitCalculator_GetBestSaleMethod(t *testing.T) {
	type args struct {
		item         *Item
		listings     *[]*db.Listing
		sales        *[]*db.Sale
		playerServer *PlayerInfo
	}
	tests := []struct {
		name string
		args args
		want *SaleMethod
	}{
		{
			name: "Unmarketable item that is sellable to vendor can still make money",
			args: args{
				item: &Item{
					Id:               5,
					MarketProhibited: true,
					CanBeTraded:      true,
					ExchangeMethods: &[]ExchangeMethod{
						GilExchange{
							TokenExchange: TokenExchange{
								Value:    50,
								Quantity: 1,
							},
							NpcName: "NPC",
						},
					},
				},
				listings: nil,
				playerServer: &PlayerInfo{
					HomeServer: 1,
					DataCenter: 1,
				},
			},
			want: &SaleMethod{
				ExchangeType: "Sell to NPC",
				Value:        50,
				Quantity:     1,
				ValuePer:     50,
			},
		},
		{
			name: "Marketable item is more profitable to sell to NPC",
			args: args{
				item: &Item{
					Id:               6,
					MarketProhibited: true,
					CanBeTraded:      true,
					ExchangeMethods: &[]ExchangeMethod{
						GilExchange{
							TokenExchange: TokenExchange{
								Value:    500,
								Quantity: 1,
							},
							NpcName: "NPC",
						},
					},
				},
				listings: &[]*db.Listing{
					{Id: 1, ItemId: 6, WorldId: 1, PricePer: 99, Quantity: 1, Total: 99},
					{Id: 2, ItemId: 6, WorldId: 1, PricePer: 100, Quantity: 1, Total: 100},
				},
				playerServer: &PlayerInfo{
					HomeServer: 1,
					DataCenter: 1,
				},
			},
			want: &SaleMethod{
				ExchangeType: "Sell to NPC",
				Value:        500,
				Quantity:     1,
				ValuePer:     500,
			},
		},
		{
			name: "Marketable item is more profitable per item to sell back on the market",
			args: args{
				item: &Item{
					Id:               5,
					MarketProhibited: true,
					CanBeTraded:      true,
					ExchangeMethods: &[]ExchangeMethod{
						GilExchange{
							TokenExchange: TokenExchange{
								Value:    400,
								Quantity: 5,
							},
							NpcName: "NPC",
						},
					},
				},
				listings: &[]*db.Listing{
					{Id: 1, ItemId: 5, WorldId: 1, PricePer: 100, Quantity: 5, Total: 500},
					{Id: 2, ItemId: 5, WorldId: 1, PricePer: 101, Quantity: 1, Total: 101},
				},
				playerServer: &PlayerInfo{
					HomeServer: 1,
					DataCenter: 1,
				},
			},
			want: &SaleMethod{
				ExchangeType: "Marketboard",
				Value:        495,
				Quantity:     5,
				ValuePer:     99,
			},
		},
		// TODO add tests for sales tracking
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				itemMap := make(map[int]*Item)

				p := NewProfitCalculator(&itemMap, db.NewMockRepository())

				if got := p.GetBestSaleMethod(
					tt.args.item,
					tt.args.listings,
					tt.args.sales,
					tt.args.playerServer,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetBestSaleMethod() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestProfitCalculator_GetCostToObtain(t *testing.T) {
	type args struct {
		item         *Item
		listings     *[]*db.Listing
		playerServer *PlayerInfo
		recipeItems  []*Item
	}
	tests := []struct {
		name string
		args args
		want *ObtainInfo
	}{
		{
			name: "Item with no obtain method or market listings safely returns nothing",
			args: args{
				item: &Item{
					Id:               1,
					MarketProhibited: true,
					ObtainMethods:    nil,
				},
				listings: nil,
				playerServer: &PlayerInfo{
					HomeServer: 1,
					DataCenter: 1,
				},
				recipeItems: nil,
			},
			want: nil,
		},
		{
			name: "Item is cheapest to buy on the marketboard",
			args: args{
				item: &Item{
					Id:               2,
					MarketProhibited: false,
					ObtainMethods: &[]ExchangeMethod{
						GilExchange{
							TokenExchange: TokenExchange{
								Value:    500,
								Quantity: 1,
							},
							NpcName: "NPC",
						},
					},
				},
				listings: &[]*db.Listing{
					{Id: 1, ItemId: 2, WorldId: 1, PricePer: 300, Quantity: 3, Total: 900},
					{Id: 2, ItemId: 2, WorldId: 2, PricePer: 100, Quantity: 1, Total: 100},
				},
				playerServer: &PlayerInfo{
					HomeServer: 1,
					DataCenter: 1,
				},
				recipeItems: nil,
			},
			want: &ObtainInfo{
				ItemsRequired: []*PurchaseInfo{
					{
						Item: &Item{
							Id:               2,
							MarketProhibited: false,
							ObtainMethods: &[]ExchangeMethod{
								GilExchange{
									TokenExchange: TokenExchange{
										Value:    500,
										Quantity: 1,
									},
									NpcName: "NPC",
								},
							},
						},
						Quantity: 1,
						Server:   2,
					},
				},
				ObtainMethod:   "Market",
				Cost:           100,
				CostPerItem:    100,
				ResultQuantity: 1,
				EffortFactor:   1.05,
			},
		},
		{
			name: "Item with recipe that simplifies down to 1 ingredient is cheaper than vendor",
			args: args{
				item: &Item{
					Id:               1,
					MarketProhibited: true,
					ObtainMethods: &[]ExchangeMethod{
						GilExchange{
							TokenExchange: TokenExchange{
								Value:    50,
								Quantity: 1,
							},
							NpcName: "Expensive vendor",
						},
					},
					CanBeHq: true,
					CraftingRecipes: &[]RecipeInfo{
						{
							Yield:       1,
							CraftType:   "CARPENTER",
							RecipeLevel: 90,
							RecipeIngredients: []RecipeIngredients{
								{
									ItemId:   2,
									Quantity: 1,
								},
								{
									ItemId:   3,
									Quantity: 2,
								},
							},
						},
					},
				},
				listings: nil,
				playerServer: &PlayerInfo{
					HomeServer: 1,
					DataCenter: 1,
				},
				recipeItems: []*Item{
					{
						Id: 2,
						CraftingRecipes: &[]RecipeInfo{
							{
								Yield:       1,
								CraftType:   "ALCHEMIST",
								RecipeLevel: 90,
								RecipeIngredients: []RecipeIngredients{
									{
										ItemId:   3,
										Quantity: 4,
									},
								},
							},
						},
					},
					{
						Id: 3,
						ObtainMethods: &[]ExchangeMethod{
							GilExchange{
								TokenExchange: TokenExchange{
									Value:    5,
									Quantity: 1,
								},
								NpcName: "NPC",
							},
						},
					},
				},
			},
			want: &ObtainInfo{
				ItemsRequired: []*PurchaseInfo{
					{
						Item: &Item{
							Id: 3,
							ObtainMethods: &[]ExchangeMethod{
								GilExchange{
									TokenExchange: TokenExchange{
										Value:    5,
										Quantity: 1,
									},
									NpcName: "NPC",
								},
							},
						},
						Quantity: 6,
						Server:   1,
					},
				},
				ObtainMethod:   "CARPENTER",
				Cost:           30,
				CostPerItem:    30,
				ResultQuantity: 1,
				EffortFactor:   1.0,
			},
		},
		{
			name: "Item with recipe that has ingredients bought from market track quantity correctly",
			args: args{
				item: &Item{
					Id:               1,
					MarketProhibited: false,
					CanBeHq:          true,
					CraftingRecipes: &[]RecipeInfo{
						{
							Yield:       1,
							CraftType:   "CULINARIAN",
							RecipeLevel: 90,
							RecipeIngredients: []RecipeIngredients{
								{
									ItemId:   2,
									Quantity: 2,
								},
								{
									ItemId:   3,
									Quantity: 2,
								},
							},
						},
					},
				},
				listings: &[]*db.Listing{
					{Id: 1, ItemId: 2, WorldId: 1, PricePer: 500, Quantity: 5, Total: 2500},
					{Id: 2, ItemId: 2, WorldId: 2, PricePer: 501, Quantity: 3, Total: 1002},
				},
				playerServer: &PlayerInfo{
					HomeServer: 1,
					DataCenter: 1,
				},
				recipeItems: []*Item{
					{
						Id: 2,
					},
					{
						Id: 3,
						ObtainMethods: &[]ExchangeMethod{
							GilExchange{
								TokenExchange: TokenExchange{
									Value:    250,
									Quantity: 1,
								},
								NpcName: "NPC",
							},
						},
					},
				},
			},
			want: &ObtainInfo{
				ItemsRequired: []*PurchaseInfo{
					{
						Item: &Item{
							Id: 2,
						},
						Quantity: 3,
						Server:   2,
					},
					{
						Item: &Item{
							Id: 3,
							ObtainMethods: &[]ExchangeMethod{
								GilExchange{
									TokenExchange: TokenExchange{
										Value:    250,
										Quantity: 1,
									},
									NpcName: "NPC",
								},
							},
						},
						Quantity: 2,
						Server:   1,
					},
				},
				ObtainMethod:   "CULINARIAN",
				Cost:           1502,
				CostPerItem:    1502,
				ResultQuantity: 1,
				EffortFactor:   1.0,
			},
		},
		{
			name: "Skips recipes that have unobtainable ingredients",
			args: args{
				item: &Item{
					Id:               1,
					MarketProhibited: true,
					ObtainMethods:    nil,
					CanBeHq:          true,
					CraftingRecipes: &[]RecipeInfo{
						{
							Yield:       1,
							CraftType:   "CARPENTER",
							RecipeLevel: 90,
							RecipeIngredients: []RecipeIngredients{
								{
									ItemId:   2,
									Quantity: 1,
								},
								{
									ItemId:   3,
									Quantity: 2,
								},
							},
						},
					},
				},
				listings: nil,
				playerServer: &PlayerInfo{
					HomeServer: 1,
					DataCenter: 1,
				},
				recipeItems: []*Item{
					{
						Id: 2,
						CraftingRecipes: &[]RecipeInfo{
							{
								Yield:       1,
								CraftType:   "ALCHEMIST",
								RecipeLevel: 90,
								RecipeIngredients: []RecipeIngredients{
									{
										ItemId:   3,
										Quantity: 4,
									},
								},
							},
						},
					},
					// Unobtainable item
					{
						Id:            3,
						ObtainMethods: nil,
					},
				},
			},
			want: nil,
		},
		{
			name: "Item only available from a vendor is tracked accurately",
			args: args{
				item: &Item{
					Id:               1,
					MarketProhibited: true,
					ObtainMethods: &[]ExchangeMethod{
						GilExchange{
							TokenExchange: TokenExchange{
								Value:    500,
								Quantity: 1,
							},
							NpcName: "NPC",
						},
					},
				},
				listings: nil,
				playerServer: &PlayerInfo{
					HomeServer: 1,
					DataCenter: 1,
				},
				recipeItems: nil,
			},
			want: &ObtainInfo{
				ItemsRequired: []*PurchaseInfo{
					{
						Item: &Item{
							Id:               1,
							MarketProhibited: true,
							ObtainMethods: &[]ExchangeMethod{
								GilExchange{
									TokenExchange: TokenExchange{
										Value:    500,
										Quantity: 1,
									},
									NpcName: "NPC",
								},
							},
						},
						Quantity: 1,
						Server:   1,
					},
				},
				ObtainMethod:   "Buy with Gil",
				Cost:           500,
				CostPerItem:    500,
				ResultQuantity: 1,
				EffortFactor:   0.85,
			},
		},
		{
			name: "Item only available from a GC Seals is tracked accurately",
			args: args{
				item: &Item{
					Id:               1,
					MarketProhibited: true,
					ObtainMethods: &[]ExchangeMethod{
						GcSealExchange{
							TokenExchange: TokenExchange{
								Value:    200,
								Quantity: 1,
							},
							RankRequired: 4,
						},
					},
				},
				listings: nil,
				playerServer: &PlayerInfo{
					HomeServer:       1,
					DataCenter:       1,
					GrandCompanyRank: 5,
				},
				recipeItems: nil,
			},
			want: &ObtainInfo{
				ItemsRequired: []*PurchaseInfo{
					{
						Item: &Item{
							Id:               1,
							MarketProhibited: true,
							ObtainMethods: &[]ExchangeMethod{
								GcSealExchange{
									TokenExchange: TokenExchange{
										Value:    200,
										Quantity: 1,
									},
									RankRequired: 4,
								},
							},
						},
						Quantity: 1,
						Server:   1,
					},
				},
				ObtainMethod:   "Grand Company Seal",
				Cost:           200,
				CostPerItem:    200,
				ResultQuantity: 1,
				EffortFactor:   0.9,
			},
		},
		{
			name: "Correctly identifies cheapest obtain method from multiple options",
			args: args{
				item: &Item{
					Id:               1,
					MarketProhibited: false,
					ObtainMethods: &[]ExchangeMethod{
						GcSealExchange{
							TokenExchange: TokenExchange{
								Value:    200,
								Quantity: 1,
							},
							RankRequired: 4,
						},
						GilExchange{
							TokenExchange: TokenExchange{
								Value:    1000,
								Quantity: 1,
							},
							NpcName: "Reasonably priced merchant",
						},
					},
					CraftingRecipes: &[]RecipeInfo{
						{
							Yield:     1,
							CraftType: "WEAVER",
							RecipeIngredients: []RecipeIngredients{
								{
									ItemId:   2,
									Quantity: 1,
								},
								{
									ItemId:   3,
									Quantity: 1,
								},
							},
						},
					},
				},
				listings: &[]*db.Listing{
					{Id: 1, ItemId: 1, WorldId: 1, PricePer: 300, Quantity: 3, Total: 900},
					{Id: 2, ItemId: 1, WorldId: 2, PricePer: 100, Quantity: 1, Total: 100},
					{Id: 3, ItemId: 2, WorldId: 3, PricePer: 50, Quantity: 1, Total: 50},
					{Id: 4, ItemId: 2, WorldId: 3, PricePer: 49, Quantity: 2, Total: 98},
					{Id: 5, ItemId: 2, WorldId: 1, PricePer: 60, Quantity: 2, Total: 120},
					{Id: 6, ItemId: 3, WorldId: 1, PricePer: 20, Quantity: 1, Total: 20},
				},
				playerServer: &PlayerInfo{
					HomeServer:       1,
					DataCenter:       1,
					GrandCompanyRank: 3,
				},
				recipeItems: []*Item{
					{
						Id:               1,
						MarketProhibited: false,
						ObtainMethods: &[]ExchangeMethod{
							GcSealExchange{
								TokenExchange: TokenExchange{
									Value:    200,
									Quantity: 1,
								},
								RankRequired: 4,
							},
							GilExchange{
								TokenExchange: TokenExchange{
									Value:    1000,
									Quantity: 1,
								},
								NpcName: "Reasonably priced merchant",
							},
						},
						CraftingRecipes: &[]RecipeInfo{
							{
								Yield:     1,
								CraftType: "WEAVER",
								RecipeIngredients: []RecipeIngredients{
									{
										ItemId:   2,
										Quantity: 1,
									},
									{
										ItemId:   3,
										Quantity: 1,
									},
								},
							},
						},
					},
					{
						Id: 2,
						ObtainMethods: &[]ExchangeMethod{
							GilExchange{
								TokenExchange: TokenExchange{
									Value:    2500,
									Quantity: 1,
								},
								NpcName: "Very expensive merchant",
							},
						},
					},
					{
						Id: 3,
					},
				},
			},
			want: &ObtainInfo{
				ItemsRequired: []*PurchaseInfo{
					{
						Item: &Item{
							Id: 2,
							ObtainMethods: &[]ExchangeMethod{
								GilExchange{
									TokenExchange: TokenExchange{
										Value:    2500,
										Quantity: 1,
									},
									NpcName: "Very expensive merchant",
								},
							},
						},
						Quantity: 1,
						Server:   3,
						BuyFrom:  "",
					},
					{
						Item: &Item{
							Id: 3,
						},
						Quantity: 1,
						Server:   1,
						BuyFrom:  "",
					},
				},
				ObtainMethod:   "WEAVER",
				Cost:           70,
				CostPerItem:    70,
				ResultQuantity: 1,
				EffortFactor:   0.5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				itemMap := make(map[int]*Item, len(tt.args.recipeItems))
				for _, item := range tt.args.recipeItems {
					itemMap[item.Id] = item
				}

				repo := db.NewMockRepository()

				if tt.args.listings != nil {
					for _, listing := range *tt.args.listings {
						_, err := repo.CreateListing(*listing)

						if err != nil {
							t.Errorf("Error creating mock listing in GetCostToObtain: %v", err)
						}
					}
				}

				p := NewProfitCalculator(&itemMap, repo)

				if got := p.GetCostToObtain(
					tt.args.item,
					1,
					tt.args.listings,
					tt.args.playerServer,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetCostToObtain() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_combinePurchaseInfo(t *testing.T) {
	type args struct {
		slice1 []*PurchaseInfo
		slice2 []*PurchaseInfo
	}
	tests := []struct {
		name string
		args args
		want []*PurchaseInfo
	}{
		{
			name: "Combines two of the same item together",
			args: args{
				slice1: []*PurchaseInfo{
					{
						Item: &Item{
							Id: 1,
						},
						Quantity: 2,
						Server:   1,
					},
				},
				slice2: []*PurchaseInfo{
					{
						Item: &Item{
							Id: 1,
						},
						Quantity: 3,
						Server:   1,
					},
				},
			},
			want: []*PurchaseInfo{
				{
					Item: &Item{
						Id: 1,
					},
					Quantity: 5,
					Server:   1,
				},
			},
		},
		{
			name: "Combines same items and keeps different ones",
			args: args{
				slice1: []*PurchaseInfo{
					{
						Item: &Item{
							Id: 1,
						},
						Quantity: 5,
						Server:   1,
					},
					{
						Item: &Item{
							Id: 2,
						},
						Quantity: 2,
						Server:   2,
					},
				},
				slice2: []*PurchaseInfo{
					{
						Item: &Item{
							Id: 1,
						},
						Quantity: 5,
						Server:   1,
					},
					{
						Item: &Item{
							Id: 3,
						},
						Quantity: 1,
						Server:   2,
					},
				},
			},
			want: []*PurchaseInfo{
				{
					Item: &Item{
						Id: 1,
					},
					Quantity: 10,
					Server:   1,
				},
				{
					Item: &Item{
						Id: 2,
					},
					Quantity: 2,
					Server:   2,
				},
				{
					Item: &Item{
						Id: 3,
					},
					Quantity: 1,
					Server:   2,
				},
			},
		},
		{
			name: "Handles cases where both arrays are empty",
			args: args{
				slice1: []*PurchaseInfo{},
				slice2: []*PurchaseInfo{},
			},
			want: []*PurchaseInfo{},
		},
		{
			name: "Handles cases where one array is emtpy",
			args: args{
				slice1: []*PurchaseInfo{
					{
						Item: &Item{
							Id: 1,
						},
						Quantity: 2,
						Server:   1,
					},
				},
				slice2: []*PurchaseInfo{},
			},
			want: []*PurchaseInfo{
				{
					Item: &Item{
						Id: 1,
					},
					Quantity: 2,
					Server:   1,
				},
			},
		},
		{
			name: "Doesn't combine the same item bought on different servers",
			args: args{
				slice1: []*PurchaseInfo{
					{
						Item: &Item{
							Id: 1,
						},
						Quantity: 2,
						Server:   1,
					},
				},
				slice2: []*PurchaseInfo{
					{
						Item: &Item{
							Id: 1,
						},
						Quantity: 3,
						Server:   2,
					},
				},
			},
			want: []*PurchaseInfo{
				{
					Item: &Item{
						Id: 1,
					},
					Quantity: 2,
					Server:   1,
				},
				{
					Item: &Item{
						Id: 1,
					},
					Quantity: 3,
					Server:   2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := combinePurchaseInfo(tt.args.slice1, tt.args.slice2); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("combinePurchaseInfo() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestProfitCalculator_getIngredientsForRecipe(t *testing.T) {
	type args struct {
		numRequired  int
		recipe       *RecipeInfo
		recipeItems  []*Item
		skipCrystals bool
	}
	tests := []struct {
		name string
		args args
		want *map[int]int
	}{
		{
			name: "returns an empty map if no ingredients",
			args: args{
				numRequired: 1,
				recipe: &RecipeInfo{
					Yield:             1,
					CraftType:         "CARPENTER",
					RecipeLevel:       90,
					RecipeIngredients: []RecipeIngredients{},
				},
				recipeItems: nil,
			},
			want: &map[int]int{},
		},
		{
			name: "returns simple recipes correctly",
			args: args{
				numRequired: 1,
				recipe: &RecipeInfo{
					Yield: 1,
					RecipeIngredients: []RecipeIngredients{
						{
							ItemId:   7,
							Quantity: 2,
						},
						{
							ItemId:   8,
							Quantity: 1,
						},
					},
				},
				recipeItems: []*Item{
					{
						Id:              7,
						CraftingRecipes: nil,
					},
					{
						Id:              8,
						CraftingRecipes: nil,
					},
				},
			},
			want: &map[int]int{
				7: 2,
				8: 1,
			},
		},
		{
			name: "returns complex recipes correctly",
			args: args{
				numRequired: 1,
				recipe: &RecipeInfo{
					Yield: 1,
					RecipeIngredients: []RecipeIngredients{
						{
							ItemId:   9,
							Quantity: 10,
						},
						{
							ItemId:   10,
							Quantity: 2,
						},
					},
				},
				recipeItems: []*Item{
					{
						Id:              9,
						CraftingRecipes: nil,
					},
					{
						Id: 10,
						CraftingRecipes: &[]RecipeInfo{
							{
								Yield: 2,
								RecipeIngredients: []RecipeIngredients{
									{
										ItemId:   11,
										Quantity: 1,
									},
								},
							},
						},
					},
					{
						Id:              11,
						CraftingRecipes: nil,
					},
				},
			},
			want: &map[int]int{
				9:  10,
				11: 1,
			},
		},
		{
			name: "adds quantity to the same item as expected",
			args: args{
				numRequired: 1,
				recipe: &RecipeInfo{
					Yield: 1,
					RecipeIngredients: []RecipeIngredients{
						{
							ItemId:   11,
							Quantity: 2,
						},
						{
							ItemId:   12,
							Quantity: 1,
						},
					},
				},
				recipeItems: []*Item{
					{
						Id:              11,
						CraftingRecipes: nil,
					},
					{
						Id: 12,
						CraftingRecipes: &[]RecipeInfo{
							{
								Yield: 2,
								RecipeIngredients: []RecipeIngredients{
									{
										ItemId:   11,
										Quantity: 5,
									},
								},
							},
						},
					},
				},
			},
			want: &map[int]int{
				11: 7,
			},
		},
		{
			name: "Emulates yarn crafting correctly",
			args: args{
				numRequired: 50,
				recipe: &RecipeInfo{
					Yield: 1,
					RecipeIngredients: []RecipeIngredients{
						{
							ItemId:   20,
							Quantity: 3,
						},
						{
							ItemId:   8,
							Quantity: 4,
						},
					},
				},
				recipeItems: []*Item{
					{
						Id:              8,
						CraftingRecipes: nil,
					},
					{
						Id: 20,
						CraftingRecipes: &[]RecipeInfo{
							{
								Yield: 3,
								RecipeIngredients: []RecipeIngredients{
									{
										ItemId:   19,
										Quantity: 4,
									},
									{
										ItemId:   8,
										Quantity: 3,
									},
								},
							},
						},
					},
					{
						Id:              19,
						CraftingRecipes: nil,
					},
				},
			},
			want: &map[int]int{
				8:  350,
				19: 200,
			},
		},
		// TODO add tests including the skipping of crystals
	}
	for _, tt := range tests {

		t.Run(
			tt.name, func(t *testing.T) {

				itemMap := make(map[int]*Item, len(tt.args.recipeItems))
				if tt.args.recipeItems != nil {
					for _, item := range tt.args.recipeItems {
						itemMap[item.Id] = item
					}
				}

				repo := db.NewMockRepository()
				p := NewProfitCalculator(&itemMap, repo)

				if got := p.getIngredientsForRecipe(
					nil,
					tt.args.numRequired,
					tt.args.recipe,
					tt.args.skipCrystals,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("getIngredientsForRecipe() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}