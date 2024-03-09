package profitCalc

import (
	"reflect"
	"sort"
	"testing"
)

func TestShoppingCart_mergeWith(t *testing.T) {
	tests := []struct {
		name  string
		got   ShoppingCart
		other ShoppingCart
		want  ShoppingCart
	}{
		{
			name:  "can merge two empty carts",
			got:   ShoppingCart{},
			other: ShoppingCart{},
			want: ShoppingCart{
				ItemsToBuy:    []ShoppingItem{},
				itemsRequired: nil,
			},
		},
		{
			name: "merges local item cart with empty shopping cart successfully",
			got: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					LocalItem{
						ItemId:       1,
						Quantity:     10,
						ObtainedFrom: "NPC",
						CostPer:      5,
					},
				},
				itemsRequired: map[int]int{
					1: 10,
				},
			},
			other: ShoppingCart{},
			want: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					LocalItem{
						ItemId:       1,
						Quantity:     10,
						ObtainedFrom: "NPC",
						CostPer:      5,
					},
				},
				itemsRequired: map[int]int{
					1: 10,
				},
			},
		},
		{
			name: "merges listing item cart with empty shopping cart successfully",
			got: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					ShoppingListing{
						ItemId:       5,
						Quantity:     11,
						RetainerName: "Bob",
						listingId:    "123-123-123",
						worldId:      3,
						CostPer:      50,
					},
				},
				itemsRequired: map[int]int{
					5: 10,
				},
			},
			other: ShoppingCart{},
			want: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					ShoppingListing{
						ItemId:       5,
						Quantity:     11,
						RetainerName: "Bob",
						listingId:    "123-123-123",
						worldId:      3,
						CostPer:      50,
					},
				},
				itemsRequired: map[int]int{
					5: 10,
				},
			},
		},
		{
			name: "combines 2 carts with the same local items correctly",
			got: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					LocalItem{
						ItemId:       1,
						Quantity:     10,
						ObtainedFrom: "NPC",
						CostPer:      5,
					},
				},
				itemsRequired: map[int]int{
					1: 10,
				},
			},
			other: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					LocalItem{
						ItemId:       1,
						Quantity:     10,
						ObtainedFrom: "NPC",
						CostPer:      5,
					},
				},
				itemsRequired: map[int]int{
					1: 8,
				},
			},
			want: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					LocalItem{
						ItemId:       1,
						Quantity:     18,
						ObtainedFrom: "NPC",
						CostPer:      5,
					},
				},
				itemsRequired: map[int]int{
					1: 18,
				},
			},
		},
		{
			name: "combines 2 carts with the same listing items correctly",
			got: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					ShoppingListing{
						ItemId:       1,
						Quantity:     15,
						RetainerName: "Jimmy",
						listingId:    "321-321-321",
						worldId:      3,
						CostPer:      15,
					},
				},
				itemsRequired: map[int]int{
					1: 13,
				},
			},
			other: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					ShoppingListing{
						ItemId:       1,
						Quantity:     9,
						RetainerName: "Benjamin",
						listingId:    "234-678-890",
						worldId:      3,
						CostPer:      21,
					},
				},
				itemsRequired: map[int]int{
					1: 7,
				},
			},
			want: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					ShoppingListing{
						ItemId:       1,
						Quantity:     9,
						RetainerName: "Benjamin",
						listingId:    "234-678-890",
						worldId:      3,
						CostPer:      21,
					},
					ShoppingListing{
						ItemId:       1,
						Quantity:     15,
						RetainerName: "Jimmy",
						listingId:    "321-321-321",
						worldId:      3,
						CostPer:      15,
					},
				},
				itemsRequired: map[int]int{
					1: 20,
				},
			},
		},
		{
			name: "combines 2 carts with different items correctly",
			got: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					LocalItem{
						ItemId:       1,
						Quantity:     10,
						ObtainedFrom: "NPC",
						CostPer:      5,
					},
				},
				itemsRequired: map[int]int{
					1: 10,
				},
			},
			other: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					ShoppingListing{
						ItemId:       2,
						Quantity:     5,
						RetainerName: "Fred",
						listingId:    "123-456-789",
						worldId:      2,
						CostPer:      15,
					},
				},
				itemsRequired: map[int]int{
					2: 3,
				},
			},
			want: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					LocalItem{
						ItemId:       1,
						Quantity:     10,
						ObtainedFrom: "NPC",
						CostPer:      5,
					},
					ShoppingListing{
						ItemId:       2,
						Quantity:     5,
						RetainerName: "Fred",
						listingId:    "123-456-789",
						worldId:      2,
						CostPer:      15,
					},
				},
				itemsRequired: map[int]int{
					1: 10,
					2: 3,
				},
			},
		},
		{
			name: "removes unecessary extras from cart",
			got: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					ShoppingListing{
						ItemId:       15,
						Quantity:     999,
						RetainerName: "Bob the Crystal Seller",
						listingId:    "987-654-321",
						worldId:      3,
						CostPer:      6,
					},
				},
				itemsRequired: map[int]int{
					15: 50,
				},
			},
			other: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					ShoppingListing{
						ItemId:       15,
						Quantity:     999,
						RetainerName: "Jim the Crystal Seller",
						listingId:    "000-000-000",
						worldId:      3,
						CostPer:      4,
					},
				},
				itemsRequired: map[int]int{
					15: 30,
				},
			},
			want: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					ShoppingListing{
						ItemId:       15,
						Quantity:     999,
						RetainerName: "Bob the Crystal Seller",
						listingId:    "987-654-321",
						worldId:      3,
						CostPer:      6,
					},
				},
				itemsRequired: map[int]int{
					15: 80,
				},
			},
		},
		{
			name: "removes unecessary items but keeps new ones",
			got: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					LocalItem{
						ItemId:       1,
						Quantity:     9,
						ObtainedFrom: "NPC",
						CostPer:      5,
					},
					ShoppingListing{
						ItemId:       15,
						Quantity:     999,
						RetainerName: "Bob the Crystal Seller",
						listingId:    "987-654-321",
						worldId:      3,
						CostPer:      6,
					},
				},
				itemsRequired: map[int]int{
					1:  8,
					15: 50,
				},
			},
			other: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					ShoppingListing{
						ItemId:       15,
						Quantity:     999,
						RetainerName: "Jim the Crystal Seller",
						listingId:    "000-000-000",
						worldId:      3,
						CostPer:      4,
					},
					LocalItem{
						ItemId:       1,
						Quantity:     5,
						ObtainedFrom: "NPC",
						CostPer:      5,
					},
				},
				itemsRequired: map[int]int{
					1:  2,
					15: 30,
				},
			},
			want: ShoppingCart{
				ItemsToBuy: []ShoppingItem{
					LocalItem{
						ItemId:       1,
						Quantity:     10,
						ObtainedFrom: "NPC",
						CostPer:      5,
					},
					ShoppingListing{
						ItemId:       15,
						Quantity:     999,
						RetainerName: "Bob the Crystal Seller",
						listingId:    "987-654-321",
						worldId:      3,
						CostPer:      6,
					},
				},
				itemsRequired: map[int]int{
					1:  10,
					15: 80,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := tt.got

				got.mergeWith(tt.other)

				sort.Slice(
					got.ItemsToBuy, func(i, j int) bool {
						return got.ItemsToBuy[i].GetHash() < got.ItemsToBuy[j].GetHash()
					},
				)

				if !reflect.DeepEqual(tt.want, got) {
					t.Errorf("mergeWith() = got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
