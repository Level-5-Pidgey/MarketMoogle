package profitCalc

import (
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
	"github.com/level-5-pidgey/MarketMoogle/db"
	"github.com/level-5-pidgey/MarketMoogle/profit/exchange"
	"reflect"
	"sort"
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
					ExchangeMethods: &[]exchange.Method{
						exchange.NewGilExchange(50, "NPC", ""),
					},
				},
				listings: nil,
				playerServer: &PlayerInfo{
					HomeServer: 1,
					DataCenter: 1,
				},
			},
			want: &SaleMethod{
				ExchangeType:      readertype.Gil,
				Value:             50,
				Quantity:          1,
				ValuePer:          50,
				SaleVelocity:      0.0001,
				CompetitionFactor: 1.0,
			},
		},
		{
			name: "Marketable item is more profitable to sell to NPC",
			args: args{
				item: &Item{
					Id:               6,
					MarketProhibited: true,
					CanBeTraded:      true,
					ExchangeMethods: &[]exchange.Method{
						exchange.NewGilExchange(500, "NPC", ""),
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
				ExchangeType:      readertype.Gil,
				Value:             500,
				Quantity:          1,
				ValuePer:          500,
				SaleVelocity:      0.0001,
				CompetitionFactor: 0.5,
			},
		},
		{
			name: "Marketable item is more profitable per item to sell back on the market",
			args: args{
				item: &Item{
					Id:               5,
					MarketProhibited: true,
					CanBeTraded:      true,
					ExchangeMethods: &[]exchange.Method{
						exchange.NewGilExchange(80, "NPC", ""),
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
				ExchangeType:      readertype.Marketboard,
				Value:             495,
				Quantity:          5,
				ValuePer:          99,
				SaleVelocity:      0.0001,
				CompetitionFactor: 0.5,
			},
		},
		// TODO add tests for sales tracking
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				itemMap := make(map[int]*Item)

				p := NewProfitCalculator(
					&itemMap,
					nil,
					nil,
					db.NewMockRepository(),
				)

				if got := p.GetBestSaleMethod(
					tt.args.item,
					tt.args.listings,
					tt.args.sales,
					tt.args.playerServer,
					false,
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
		want *ObtainMethod
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
					ObtainMethods: &[]exchange.Method{
						exchange.NewGilExchange(500, "NPC", ""),
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
			want: &ObtainMethod{
				ShoppingCart: ShoppingCart{
					ItemsToBuy: []ShoppingItem{
						ShoppingListing{
							ItemId:    2,
							Quantity:  1,
							listingId: 2,
							worldId:   2,
							CostPer:   100,
						},
					},
					itemsRequired: map[int]int{2: 1},
				},
				ObtainMethod: "Market",
				Quantity:     1,
				EffortFactor: 1.05,
			},
		},
		{
			name: "Item with recipe that simplifies down to 1 ingredient is cheaper than vendor",
			args: args{
				item: &Item{
					Id:               1,
					MarketProhibited: true,
					ObtainMethods: &[]exchange.Method{
						exchange.NewGilExchange(75, "Expensive vendor", ""),
					},
					CanBeHq: true,
					CraftingRecipes: &[]RecipeInfo{
						{
							Yield:       1,
							JobRequired: readertype.JobCarpenter,
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
					JobLevels: map[readertype.Job]int{
						readertype.JobCarpenter: 90,
						readertype.JobAlchemist: 90,
					},
				},
				recipeItems: []*Item{
					{
						Id: 2,
						CraftingRecipes: &[]RecipeInfo{
							{
								Yield:       1,
								JobRequired: readertype.JobAlchemist,
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
						ObtainMethods: &[]exchange.Method{
							exchange.NewGilExchange(5, "NPC", ""),
						},
					},
				},
			},
			want: &ObtainMethod{
				ShoppingCart: ShoppingCart{
					ItemsToBuy: []ShoppingItem{
						LocalItem{
							ItemId:       3,
							Quantity:     6,
							ObtainedFrom: "Buy from NPC",
							CostPer:      5,
						},
					},
					itemsRequired: map[int]int{3: 6},
				},
				ObtainMethod: "Craft with Carpenter",
				Quantity:     1,
				EffortFactor: 1.0,
			},
		},
		{
			name: "Item with recipe that has ingredients bought from market track Quantity correctly",
			args: args{
				item: &Item{
					Id:               1,
					MarketProhibited: false,
					CanBeHq:          true,
					CraftingRecipes: &[]RecipeInfo{
						{
							Yield:       1,
							JobRequired: readertype.JobCulinarian,
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
					JobLevels: map[readertype.Job]int{
						readertype.JobCulinarian: 90,
					},
				},
				recipeItems: []*Item{
					{
						Id: 2,
					},
					{
						Id: 3,
						ObtainMethods: &[]exchange.Method{
							exchange.NewGilExchange(250, "NPC", ""),
						},
					},
				},
			},
			want: &ObtainMethod{
				ShoppingCart: ShoppingCart{
					ItemsToBuy: []ShoppingItem{
						ShoppingListing{
							ItemId:    2,
							Quantity:  3,
							listingId: 2,
							worldId:   2,
							CostPer:   501,
						},
						LocalItem{
							ItemId:       3,
							Quantity:     2,
							ObtainedFrom: "Buy from NPC",
							CostPer:      250,
						},
					},
					itemsRequired: map[int]int{
						2: 2,
						3: 2,
					},
				},
				ObtainMethod: "Craft with Culinarian",
				Quantity:     1,
				EffortFactor: 1.0,
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
							JobRequired: readertype.JobCarpenter,
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
					JobLevels: map[readertype.Job]int{
						readertype.JobCarpenter: 90,
						readertype.JobAlchemist: 90,
					},
				},
				recipeItems: []*Item{
					{
						Id: 2,
						CraftingRecipes: &[]RecipeInfo{
							{
								Yield:       1,
								JobRequired: readertype.JobAlchemist,
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
					ObtainMethods: &[]exchange.Method{
						exchange.NewGilExchange(500, "NPC", ""),
					},
				},
				listings: nil,
				playerServer: &PlayerInfo{
					HomeServer: 1,
					DataCenter: 1,
				},
				recipeItems: nil,
			},
			want: &ObtainMethod{
				ShoppingCart: ShoppingCart{
					ItemsToBuy: []ShoppingItem{
						LocalItem{
							ItemId:       1,
							Quantity:     1,
							ObtainedFrom: "Buy from NPC",
							CostPer:      500,
						},
					},
					itemsRequired: map[int]int{
						1: 1,
					},
				},
				ObtainMethod: "Gil",
				Quantity:     1,
				EffortFactor: 0.85,
			},
		},
		{
			name: "Item only available from a GC Seals is tracked accurately",
			args: args{
				item: &Item{
					Id:               1,
					MarketProhibited: true,
					ObtainMethods: &[]exchange.Method{
						exchange.NewGcSealExchange(200, "", "", readertype.Corporal),
					},
				},
				listings: nil,
				playerServer: &PlayerInfo{
					HomeServer:       1,
					DataCenter:       1,
					GrandCompanyRank: readertype.SergeantThirdClass,
				},
				recipeItems: nil,
			},
			want: &ObtainMethod{
				ShoppingCart: ShoppingCart{
					ItemsToBuy: []ShoppingItem{
						LocalItem{
							ItemId:       1,
							Quantity:     1,
							ObtainedFrom: "Exchange Grand Company Seals (Rank: Corporal)",
							CostPer:      500, // This is the default gil cost for a currency grind
						},
					},
					itemsRequired: map[int]int{
						1: 1,
					},
				},
				ObtainMethod: "Grand Company Seal",
				Quantity:     1,
				EffortFactor: 0.9,
			},
		},
		{
			name: "Correctly identifies cheapest obtain method from multiple options",
			args: args{
				item: &Item{
					Id:               1,
					MarketProhibited: false,
					ObtainMethods: &[]exchange.Method{
						exchange.NewGcSealExchange(200, "", "", readertype.PrivateSecondClass),
						exchange.NewGilExchange(1000, "Reasonably priced merchant", ""),
					},
					CraftingRecipes: &[]RecipeInfo{
						{
							Yield:       1,
							JobRequired: readertype.JobWeaver,
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
					GrandCompanyRank: readertype.PrivateThirdClass,
					JobLevels: map[readertype.Job]int{
						readertype.JobWeaver: 90,
					},
				},
				recipeItems: []*Item{
					{
						Id:               1,
						MarketProhibited: false,
						ObtainMethods: &[]exchange.Method{
							exchange.NewGcSealExchange(200, "", "", readertype.PrivateSecondClass),
							exchange.NewGilExchange(1000, "Reasonably priced merchant", ""),
						},
						CraftingRecipes: &[]RecipeInfo{
							{
								Yield:       1,
								JobRequired: readertype.JobWeaver,
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
						ObtainMethods: &[]exchange.Method{
							exchange.NewGilExchange(2500, "Very expensive merchant", ""),
						},
					},
					{
						Id: 3,
					},
				},
			},
			want: &ObtainMethod{
				ShoppingCart: ShoppingCart{
					ItemsToBuy: []ShoppingItem{
						ShoppingListing{
							ItemId:    2,
							Quantity:  1,
							listingId: 3,
							worldId:   3,
							CostPer:   50,
						},
						ShoppingListing{
							ItemId:    3,
							Quantity:  1,
							listingId: 6,
							worldId:   1,
							CostPer:   20,
						},
					},
					itemsRequired: map[int]int{
						2: 1,
						3: 1,
					},
				},
				ObtainMethod: "Craft with Weaver",
				Quantity:     1,
				EffortFactor: 0.95,
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
							t.Errorf("Error creating mock listing in GetCheapestObtainMethod: %v", err)
						}
					}
				}

				p := NewProfitCalculator(
					&itemMap,
					nil,
					nil,
					repo,
				)

				got := p.GetCheapestObtainMethod(
					tt.args.item,
					1,
					tt.args.listings,
					tt.args.playerServer,
				)

				if got != nil && got.ShoppingCart.ItemsToBuy != nil {
					sort.Slice(
						got.ShoppingCart.ItemsToBuy, func(i, j int) bool {
							return got.ShoppingCart.ItemsToBuy[i].GetHash() < got.ShoppingCart.ItemsToBuy[j].GetHash()
						},
					)
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetCheapestObtainMethod() = %v, want %v", got, tt.want)
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
						ItemId:   1,
						Quantity: 2,
						Server:   1,
					},
				},
				slice2: []*PurchaseInfo{
					{
						ItemId:   1,
						Quantity: 3,
						Server:   1,
					},
				},
			},
			want: []*PurchaseInfo{
				{
					ItemId:   1,
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
						ItemId:   1,
						Quantity: 5,
						Server:   1,
					},
					{
						ItemId:   2,
						Quantity: 2,
						Server:   2,
					},
				},
				slice2: []*PurchaseInfo{
					{
						ItemId:   1,
						Quantity: 5,
						Server:   1,
					},
					{
						ItemId:   3,
						Quantity: 1,
						Server:   2,
					},
				},
			},
			want: []*PurchaseInfo{
				{
					ItemId:   1,
					Quantity: 10,
					Server:   1,
				},
				{
					ItemId:   2,
					Quantity: 2,
					Server:   2,
				},
				{
					ItemId:   3,
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
						ItemId:   1,
						Quantity: 2,
						Server:   1,
					},
				},
				slice2: []*PurchaseInfo{},
			},
			want: []*PurchaseInfo{
				{
					ItemId:   1,
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
						ItemId:   1,
						Quantity: 2,
						Server:   1,
					},
				},
				slice2: []*PurchaseInfo{
					{
						ItemId:   1,
						Quantity: 3,
						Server:   2,
					},
				},
			},
			want: []*PurchaseInfo{
				{
					ItemId:   1,
					Quantity: 2,
					Server:   1,
				},
				{
					ItemId:   1,
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
					JobRequired:       readertype.JobCarpenter,
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
			name: "adds Quantity to the same item as expected",
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
				p := NewProfitCalculator(
					&itemMap,
					nil,
					nil,
					repo,
				)

				if got := p.getPossibleSubItems(
					nil,
					tt.args.numRequired,
					tt.args.recipe,
					tt.args.skipCrystals,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("getPossibleSubItems() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
