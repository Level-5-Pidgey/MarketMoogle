package profitCalc

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
)

type ShoppingListing struct {
	ItemId       int
	Quantity     int
	RetainerName string
	listingId    int
	worldId      int
	CostPer      int
}

func (s ShoppingListing) GetTotalCost() int {
	return s.CostPer * s.Quantity
}

func (s ShoppingListing) GetCostPer() int {
	return s.CostPer
}

func (s ShoppingListing) GetItemId() int {
	return s.ItemId
}

func (s ShoppingListing) GetQuantity() int {
	return s.Quantity
}

func (s ShoppingListing) BuyFrom() string {
	return s.RetainerName
}

func (s ShoppingListing) GetHash() string {
	return fmt.Sprintf("%d_%d", s.ItemId, s.listingId)
}

type LocalItem struct {
	ItemId       int
	Quantity     int
	ObtainedFrom string
	CostPer      float64
}

func (l LocalItem) GetTotalCost() int {
	return l.GetCostPer() * l.Quantity
}

func (l LocalItem) GetCostPer() int {
	return int(math.Ceil(l.CostPer))
}

func (l LocalItem) GetItemId() int {
	return l.ItemId
}

func (l LocalItem) GetQuantity() int {
	return l.Quantity
}

func (l LocalItem) BuyFrom() string {
	return l.ObtainedFrom
}

func (l LocalItem) GetHash() string {
	return strconv.Itoa(l.ItemId)
}

type ShoppingItem interface {
	GetItemId() int
	GetTotalCost() int
	GetCostPer() int
	GetQuantity() int
	BuyFrom() string
	GetHash() string
}

type ShoppingCart struct {
	ItemsToBuy []ShoppingItem

	itemsRequired map[int]int
}

func (currentCart *ShoppingCart) MarshalJSON() ([]byte, error) {
	sort.Slice(
		currentCart.ItemsToBuy, func(i, j int) bool {
			return currentCart.ItemsToBuy[i].GetItemId() < currentCart.ItemsToBuy[j].GetItemId()
		},
	)

	return json.Marshal(currentCart.ItemsToBuy)
}

func (currentCart *ShoppingCart) mergeWith(other ShoppingCart) {
	for itemId, quantity := range other.itemsRequired {
		currentCart.itemsRequired[itemId] += quantity
	}

	combinedLength := len(currentCart.ItemsToBuy) + len(other.ItemsToBuy)
	itemQuantitiesById := make(map[int]int, combinedLength)
	hashedItemsToBuy := make(map[string]ShoppingItem, combinedLength)

	for _, item := range currentCart.ItemsToBuy {
		itemQuantitiesById[item.GetItemId()] += item.GetQuantity()
		hashedItemsToBuy[item.GetHash()] = item
	}

	for _, item := range other.ItemsToBuy {
		itemId := item.GetItemId()
		existingQuantity, exists := itemQuantitiesById[itemId]
		requiredQuantity := currentCart.itemsRequired[itemId]

		if !exists || existingQuantity < currentCart.itemsRequired[itemId] {
			itemQuantitiesById[itemId] += item.GetQuantity()

			if localItem, ok := item.(LocalItem); ok {
				localItem.Quantity = requiredQuantity
				hashedItemsToBuy[localItem.GetHash()] = localItem
			} else {
				hashedItemsToBuy[item.GetHash()] = item
			}
		}
	}

	newList := make([]ShoppingItem, 0, len(hashedItemsToBuy))
	for _, item := range hashedItemsToBuy {
		newList = append(newList, item)
	}

	currentCart.ItemsToBuy = newList
}
