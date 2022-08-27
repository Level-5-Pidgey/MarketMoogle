/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (profitprovider.go) is part of MarketMoogleAPI and is released under the GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package providers

import (
	"MarketMoogleAPI/business/database"
	schema "MarketMoogleAPI/core/graph/model"
	"context"
	"log"
	"math"
)

type ItemProfitProvider struct {
	maxValue       int
	returnUnlisted bool

	//Relevant providers
	recipeProvider      database.RecipeProvider
	marketboardProvider database.MarketBoardProvider
	itemProvider        database.ItemProvider
}

func NewItemProfitProvider(recipeProv database.RecipeProvider, mbProv database.MarketBoardProvider, itemProv database.ItemProvider) *ItemProfitProvider {
	prov := ItemProfitProvider{
		maxValue:            math.MaxInt32,
		returnUnlisted:      true,
		recipeProvider:      recipeProv,
		marketboardProvider: mbProv,
		itemProvider:        itemProv,
	}

	return &prov
}

func (profitProv ItemProfitProvider) GetRecipePurchaseInfo(componentItem *schema.Item, mbEntry *schema.MarketboardEntry, homeServer string, buyFromOtherServers *bool, count int) *schema.RecipePurchaseInfo {
	result := schema.RecipePurchaseInfo{
		Item:            componentItem,
		ServerToBuyFrom: homeServer,
		BuyFromVendor:   false,
		SingleCost:      profitProv.maxValue,
		TotalCost:       profitProv.maxValue,
		Quantity:        count,
	}

	//Always buy from a vendor if possible, before going to the market board
	if componentItem.BuyFromVendorValue != nil {
		result.SingleCost = *componentItem.BuyFromVendorValue
		result.TotalCost = *componentItem.BuyFromVendorValue * result.Quantity
		result.BuyFromVendor = true
	} else {
		if mbEntry == nil {
			return &result
		}

		cheapestEntry := &schema.MarketEntry{}

		//Get cheapest value for the item (on own server or other server if allowed)
		if buyFromOtherServers != nil && !*buyFromOtherServers {
			cheapestEntry = profitProv.GetCheapestOnServer(mbEntry, homeServer)
		} else if mbEntry.CurrentMinPrice != nil {
			cheapestEntry = profitProv.GetCheapestOnDc(mbEntry)
		}

		//If there are no market entries available, grab the listed average price instead
		if cheapestEntry.TotalPrice == profitProv.maxValue {
			result.SingleCost = int(mbEntry.CurrentAveragePrice)
		}

		result.SingleCost = cheapestEntry.PricePer
		result.TotalCost = cheapestEntry.PricePer * result.Quantity
		result.ServerToBuyFrom = cheapestEntry.Server
	}

	return &result
}

func (profitProv ItemProfitProvider) GetVendorFlipProfit(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string) (int, error) {
	if obj == nil {
		return 0, nil
	}

	marketEntry, err := profitProv.marketboardProvider.FindItemEntryOnDc(ctx, obj.ItemID, dataCenter)

	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	if marketEntry == nil {
		return 0, nil
	}

	vendorValue := 0
	if obj.BuyFromVendorValue != nil {
		vendorValue = *obj.BuyFromVendorValue
	} else {
		return 0, nil
	}

	minHomeValue := profitProv.maxValue
	for _, entry := range marketEntry.MarketEntries {
		if entry.Server == homeServer && entry.PricePer < minHomeValue {
			minHomeValue = entry.PricePer
		}
	}

	if minHomeValue == profitProv.maxValue {
		return 0, nil
	}

	return minHomeValue - vendorValue, nil
}

func (profitProv ItemProfitProvider) getItemValue(marketEntries *schema.MarketboardEntry, server string) int {
	result := 0

	cheapestOnServer := profitProv.GetCheapestOnServer(marketEntries, server)
	if cheapestOnServer.TotalPrice == profitProv.maxValue {
		result = int(marketEntries.CurrentAveragePrice)
	} else {
		result = cheapestOnServer.PricePer
	}

	return result
}

func (profitProv ItemProfitProvider) getRecipeResaleInfo(ctx context.Context, recipe *schema.Recipe, buyCrystals *bool, buyFromOtherServers *bool, homeServer string, dataCenter string) (*schema.ResaleInfo, error) {
	result := &schema.ResaleInfo{
		Profit:          0,
		ItemID:          recipe.ItemResultID,
		Quantity:        recipe.ResultQuantity,
		SingleCost:      profitProv.maxValue,
		TotalCost:       profitProv.maxValue,
		ItemsToPurchase: []*schema.RecipePurchaseInfo{},
	}

	recipeComponents := recipe.RecipeItems
	if len(recipeComponents) > 0 {
		recipeCost := 0

		for _, recipeComponent := range recipeComponents {
			if (buyCrystals == nil || !*buyCrystals) && recipeComponent.ItemID <= 19 { //19 is the highest item ID of elemental crystals
				continue
			}

			componentItem, err := profitProv.itemProvider.FindItemByItemId(ctx, recipeComponent.ItemID)
			if err != nil {
				log.Fatal(err)
				return result, err
			}

			itemMarketboardEntry, err := profitProv.marketboardProvider.FindItemEntryOnDc(ctx, componentItem.ItemID, dataCenter)
			if err != nil {
				log.Fatal(err)
				return result, err
			}

			itemPurchaseInfo := profitProv.GetRecipePurchaseInfo(componentItem, itemMarketboardEntry, homeServer, buyFromOtherServers, recipeComponent.Count)

			if itemPurchaseInfo.SingleCost != profitProv.maxValue {
				result.ItemsToPurchase = append(result.ItemsToPurchase, itemPurchaseInfo)

				//Update the cost of this recipe.
				recipeCost += itemPurchaseInfo.SingleCost * recipeComponent.Count
			} else {
				break
			}
		}

		result.TotalCost = recipeCost
		result.SingleCost = recipeCost / recipe.ResultQuantity
	}

	return result, nil
}

func (profitProv ItemProfitProvider) GetResaleInfoForItem(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string, buyCrystals *bool, buyFromOtherServers *bool) (*schema.RecipeResaleInfo, error) {
	var result = schema.RecipeResaleInfo{
		ResaleInfo: nil,
		CraftLevel: 0,
		CraftType:  "",
	}

	if obj == nil {
		return nil, nil
	}

	itemRecipes, err := profitProv.recipeProvider.FindRecipesByItemId(ctx, obj.ItemID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	//If the item has no recipes, there's no crafting profit to be had
	if len(itemRecipes) == 0 {
		return nil, nil
	}

	//Get the sell value of the item on the player's home server
	mbEntry, err := profitProv.marketboardProvider.FindItemEntryOnDc(ctx, obj.ItemID, dataCenter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if mbEntry == nil || len(mbEntry.MarketEntries) == 0 {
		return &result, nil
	}

	//Get the value of the item on the player's server (or the average going price currently)
	itemValue := profitProv.getItemValue(mbEntry, homeServer)

	craftType := schema.CrafterType("")
	craftLevel := 0

	for i, itemRecipe := range itemRecipes {
		if i != 0 {
			continue
		}

		if itemRecipe == nil {
			continue
		}

		craftType = itemRecipe.CraftedBy
		craftLevel = *itemRecipe.RecipeLevel

		recipeResaleInfo, err := profitProv.getRecipeResaleInfo(
			ctx, itemRecipe, buyCrystals, buyFromOtherServers, homeServer, dataCenter)

		if err != nil {
			log.Fatal(err)
			return &result, err
		}

		result.ResaleInfo = recipeResaleInfo
	}

	result.ResaleInfo.Profit = itemValue - result.ResaleInfo.TotalCost
	result.CraftType = craftType
	result.CraftLevel = craftLevel

	return &result, nil
}

func (profitProv ItemProfitProvider) GetCrossDcResaleProfit(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string) (*schema.ResaleInfo, error) {
	itemToPurchase := schema.RecipePurchaseInfo{
		Item:            obj,
		ServerToBuyFrom: homeServer,
		BuyFromVendor:   false,
		SingleCost:      profitProv.maxValue,
		TotalCost:       profitProv.maxValue,
		Quantity:        1,
	}

	result := schema.ResaleInfo{
		Profit:          0,
		ItemID:          0,
		Quantity:        1,
		SingleCost:      profitProv.maxValue,
		TotalCost:       profitProv.maxValue,
		ItemsToPurchase: []*schema.RecipePurchaseInfo{},
	}

	//If no item has been passed through, return early
	if obj == nil {
		return nil, nil
	}

	//Initialize Provider and find item
	marketEntry, err := profitProv.marketboardProvider.FindItemEntryOnDc(ctx, obj.ItemID, dataCenter)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	//If there are no market entries for this item, return a profit of 0.
	if marketEntry == nil {
		return &result, nil
	}

	homeEntry, awayEntry := profitProv.getHomeAndAwayItems(marketEntry, homeServer)

	//Calculate profit per item
	profit := homeEntry.PricePer - awayEntry.PricePer

	//Update result
	itemToPurchase.ServerToBuyFrom = awayEntry.Server
	itemToPurchase.SingleCost = awayEntry.PricePer
	itemToPurchase.TotalCost = awayEntry.TotalPrice

	result.SingleCost = itemToPurchase.SingleCost
	result.TotalCost = itemToPurchase.TotalCost
	result.Quantity = awayEntry.Quantity
	result.Profit = profit * awayEntry.Quantity
	result.ItemsToPurchase = append(result.ItemsToPurchase, &itemToPurchase)

	return &result, nil
}

func (profitProv ItemProfitProvider) getHomeAndAwayItems(marketEntry *schema.MarketboardEntry, homeServer string) (*schema.MarketEntry, *schema.MarketEntry) {
	cheapestOnDc := profitProv.GetCheapestOnDc(marketEntry)
	homeEntry := &schema.MarketEntry{PricePer: profitProv.maxValue}
	awayEntry := &schema.MarketEntry{PricePer: profitProv.maxValue}

	//See if you can flip on their own server (buy the 2nd cheapest)
	sameServerFlip := false
	if cheapestOnDc.Server == homeServer {
		for _, entry := range marketEntry.MarketEntries {
			if entry == nil {
				continue
			}

			if entry.Server == homeServer && cheapestOnDc.TotalPrice != entry.TotalPrice {
				awayEntry = entry
				sameServerFlip = true
				break
			}
		}
	} else {
		awayEntry = cheapestOnDc
	}

	homeEntry = profitProv.GetCheapestOnServer(marketEntry, homeServer)

	//If you're flipping on your own server, return values in opposite order
	//(so higher priced item is the "home" to calculate profits properly)
	if sameServerFlip {
		return awayEntry, homeEntry
	}

	return homeEntry, awayEntry
}

func (profitProv ItemProfitProvider) GetCheapestOnServer(entry *schema.MarketboardEntry, server string) *schema.MarketEntry {
	result := &schema.MarketEntry{
		ServerID:     0,
		Server:       server,
		Quantity:     1,
		TotalPrice:   profitProv.maxValue,
		PricePer:     profitProv.maxValue,
		Hq:           false,
		IsCrafted:    false,
		RetainerName: nil,
	}

	for _, marketEntry := range entry.MarketEntries {
		if marketEntry.Server != server {
			continue
		}

		if result.PricePer > marketEntry.PricePer {
			result = marketEntry
		}
	}

	return result
}

func (profitProv ItemProfitProvider) GetCheapestOnDc(entry *schema.MarketboardEntry) *schema.MarketEntry {
	result := &schema.MarketEntry{PricePer: profitProv.maxValue}

	for _, marketEntry := range entry.MarketEntries {
		if result.PricePer > marketEntry.PricePer {
			result = marketEntry
		}
	}

	return result
}
