/*
 * Copyright (c) 2022-2023 Carl Alexander Bird.
 * This file (profitprovider.go) is part of MarketMoogleAPI and is released under the GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package providers

import (
	interfaces "MarketMoogleAPI/business/database"
	schema "MarketMoogleAPI/core/graph/model"
	"context"
	"errors"
	"log"
	"math"
)

type ItemProfitProvider struct {
	maxValue       int
	returnUnlisted bool

	//Relevant providers
	recipeProvider      interfaces.RecipeProvider
	marketboardProvider interfaces.MarketBoardProvider
	itemProvider        interfaces.ItemProvider
}

func NewItemProfitProvider(recipeProv interfaces.RecipeProvider, mbProv interfaces.MarketBoardProvider, itemProv interfaces.ItemProvider) *ItemProfitProvider {
	prov := ItemProfitProvider{
		maxValue:            math.MaxInt32,
		returnUnlisted:      true,
		recipeProvider:      recipeProv,
		marketboardProvider: mbProv,
		itemProvider:        itemProv,
	}

	return &prov
}

func (profitProv ItemProfitProvider) GetComponentCostInfo(componentItem *schema.Item, mbEntry *schema.MarketboardEntry, homeServer string, buyFromOtherServers *bool, count int) *schema.ItemCostInfo {
	result := schema.ItemCostInfo{
		Item:            componentItem,
		ServerToBuyFrom: homeServer,
		PricePer:        profitProv.maxValue,
		TotalCost:       profitProv.maxValue,
		Quantity:        count,
	}

	//Always buy from a vendor if possible, before going to the market board
	if componentItem.BuyFromVendorValue != nil {
		result.PricePer = *componentItem.BuyFromVendorValue
		result.TotalCost = *componentItem.BuyFromVendorValue * 1
		result.BuyFromVendor = true
	} else {
		if mbEntry == nil {
			return &result
		}

		cheapestEntry := &schema.MarketEntry{}

		//Get cheapest value for the item (on own server or other server if allowed)
		if buyFromOtherServers != nil && !*buyFromOtherServers {
			cheapestEntry = profitProv.GetCheapestOnServer(mbEntry, homeServer)
		} else {
			cheapestEntry = profitProv.GetCheapestOnDataCenter(mbEntry)
		}

		//If there are no market entries available, grab the listed average price instead
		if cheapestEntry.TotalCost == profitProv.maxValue {
			result.PricePer = int(mbEntry.CurrentAveragePrice)
		}

		result.PricePer = cheapestEntry.PricePer
		result.TotalCost = cheapestEntry.TotalCost
		result.ServerToBuyFrom = cheapestEntry.Server
		result.Quantity = cheapestEntry.Quantity
	}

	return &result
}

func (profitProv ItemProfitProvider) GetVendorFlipProfit(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string) (int, error) {
	if obj == nil {
		return 0, nil
	}

	marketEntry, err := profitProv.marketboardProvider.FindItemEntryAcrossDataCenter(ctx, obj.Id, dataCenter)

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
	if cheapestOnServer.TotalCost == profitProv.maxValue {
		result = int(marketEntries.CurrentAveragePrice)
	} else {
		result = cheapestOnServer.PricePer
	}

	return result
}

func (profitProv ItemProfitProvider) getRecipeProfitInfo(ctx context.Context, recipe *schema.Recipe, itemValue int, buyCrystals *bool, buyFromOtherServers *bool, homeServer string, dataCenter string) (*schema.ProfitInfo, error) {
	result := &schema.ProfitInfo{
		ItemID:          recipe.ItemResultID,
		Quantity:        recipe.ResultQuantity,
		SingleCost:      profitProv.maxValue,
		TotalCost:       profitProv.maxValue,
		ItemsToPurchase: []*schema.ItemCostInfo{},
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

			itemMarketboardEntry, err := profitProv.marketboardProvider.FindItemEntryAcrossDataCenter(ctx, componentItem.Id, dataCenter)
			if err != nil {
				log.Fatal(err)
				return result, err
			}

			itemPurchaseInfo := profitProv.GetComponentCostInfo(componentItem, itemMarketboardEntry, homeServer, buyFromOtherServers, recipeComponent.Count)

			if itemPurchaseInfo.PricePer != profitProv.maxValue {
				result.ItemsToPurchase = append(result.ItemsToPurchase, itemPurchaseInfo)

				//Update the cost of this recipe.
				recipeCost += itemPurchaseInfo.PricePer * recipeComponent.Count
			} else {
				break
			}
		}

		result.TotalCost = recipeCost
		result.SingleCost = recipeCost / recipe.ResultQuantity
		result.Profit = itemValue - recipeCost
	}

	return result, nil
}

func (profitProv ItemProfitProvider) GetRecipeProfitForItem(ctx context.Context, item *schema.Item, dataCenter string, homeServer string, buyCrystals *bool, buyFromOtherServers *bool) (*schema.RecipeProfitInfo, error) {
	var result = schema.RecipeProfitInfo{
		ResaleInfo: &schema.ProfitInfo{
			Profit:     0,
			ItemID:     item.Id,
			SingleCost: 0,
			TotalCost:  0,
		},
		CraftLevel: 0,
		CraftType:  "",
	}

	if item == nil {
		return nil, nil
	}

	itemRecipes, err := profitProv.recipeProvider.FindRecipesByItemId(ctx, item.Id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	//If the item has no recipes, there's no crafting profit to be had
	if len(itemRecipes) == 0 {
		return nil, nil
	}

	//Get the sell value of the item on the player's home server
	mbEntry, err := profitProv.marketboardProvider.FindItemEntryAcrossDataCenter(ctx, item.Id, dataCenter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	//Get the value of the item on the player's server (or the average going price currently)
	itemValue := profitProv.getItemValue(mbEntry, homeServer)

	for _, itemRecipe := range itemRecipes {
		if itemRecipe == nil {
			continue
		}

		recipeResaleInfo, err := profitProv.getRecipeProfitInfo(
			ctx, itemRecipe, itemValue, buyCrystals, buyFromOtherServers, homeServer, dataCenter)

		if err != nil {
			log.Fatal(err)
			return &result, err
		}

		if recipeResaleInfo.Profit < result.ResaleInfo.Profit {
			continue
		}

		result.ResaleInfo = recipeResaleInfo
		result.CraftType = itemRecipe.CraftedBy
		result.CraftLevel = *itemRecipe.RecipeLevel
		result.ResaleInfo.Profit = recipeResaleInfo.Profit
	}

	return &result, nil
}

func (profitProv ItemProfitProvider) GetCrossDcResaleProfit(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string) (*schema.ProfitInfo, error) {
	itemToPurchase := schema.ItemCostInfo{
		Item:            obj,
		ServerToBuyFrom: homeServer,
		BuyFromVendor:   false,
		PricePer:        profitProv.maxValue,
		TotalCost:       profitProv.maxValue,
		Quantity:        1,
	}

	result := schema.ProfitInfo{
		Profit:          0,
		ItemID:          1,
		Quantity:        1,
		SingleCost:      profitProv.maxValue,
		TotalCost:       profitProv.maxValue,
		ItemsToPurchase: []*schema.ItemCostInfo{},
	}

	//If no item has been passed through, return early
	if obj == nil {
		return nil, nil
	}

	result.ItemID = obj.Id

	//Initialize Provider and find item
	marketEntry, err := profitProv.marketboardProvider.FindItemEntryAcrossDataCenter(ctx, obj.Id, dataCenter)
	if err != nil {
		return nil, err
	}

	if marketEntry == nil {
		return nil, errors.New("item searched for is not marketable")
	}

	homeEntry, awayEntry := profitProv.getHomeAndAwayItems(marketEntry, homeServer)

	//Calculate profit per item
	profit := 0
	if homeEntry.PricePer != profitProv.maxValue {
		profit = (homeEntry.PricePer * awayEntry.Quantity) - awayEntry.TotalCost
	}

	//Update result
	itemToPurchase.ServerToBuyFrom = awayEntry.Server
	itemToPurchase.PricePer = awayEntry.PricePer
	itemToPurchase.TotalCost = awayEntry.TotalCost
	itemToPurchase.Quantity = awayEntry.Quantity

	result.SingleCost = itemToPurchase.PricePer
	result.TotalCost = itemToPurchase.TotalCost
	result.Quantity = awayEntry.Quantity
	result.Profit = profit
	result.ItemsToPurchase = append(result.ItemsToPurchase, &itemToPurchase)

	return &result, nil
}

func (profitProv ItemProfitProvider) getHomeAndAwayItems(marketEntry *schema.MarketboardEntry, homeServer string) (*schema.MarketEntry, *schema.MarketEntry) {
	homeEntry := &schema.MarketEntry{PricePer: profitProv.maxValue, TotalCost: profitProv.maxValue}
	awayEntry := &schema.MarketEntry{PricePer: profitProv.maxValue, TotalCost: profitProv.maxValue}

	//If there are no market entries for this item, return early
	if len(marketEntry.MarketEntries) == 0 {
		return homeEntry, awayEntry
	}

	homeEntry = profitProv.GetCheapestOnServer(marketEntry, homeServer)
	priceEntries := marketEntry.MarketEntries

	bestProfit := 0
	for _, entry := range priceEntries {
		profit := 0

		if entry.Server == homeServer {
			//Even though the profit margin would be higher, players will only buy the cheapest items
			if entry.PricePer >= awayEntry.PricePer {
				continue
			}

			profit = entry.PricePer - homeEntry.PricePer
		} else {
			profit = homeEntry.PricePer - entry.PricePer
		}

		if profit <= 0 {
			continue
		}

		if profit > bestProfit {
			bestProfit = profit
			awayEntry = entry
		}
	}

	//If you're flipping on your own server, return values in opposite order
	//(so higher priced item is the "home" to calculate profits properly)
	if awayEntry.Server == homeServer {
		return awayEntry, homeEntry
	}

	return homeEntry, awayEntry
}

func (profitProv ItemProfitProvider) GetCheapestOnServer(entry *schema.MarketboardEntry, server string) *schema.MarketEntry {
	result := &schema.MarketEntry{
		ServerID:     0,
		Server:       server,
		Quantity:     1,
		TotalCost:    profitProv.maxValue,
		PricePer:     profitProv.maxValue,
		Hq:           false,
		IsCrafted:    false,
		RetainerName: nil,
	}

	for _, marketEntry := range entry.MarketEntries {
		if marketEntry.Server != server {
			continue
		}

		if marketEntry.TotalCost < result.TotalCost {
			result = marketEntry
		}
	}

	return result
}

func (profitProv ItemProfitProvider) GetCheapestOnDataCenter(entry *schema.MarketboardEntry) *schema.MarketEntry {
	result := &schema.MarketEntry{TotalCost: profitProv.maxValue}

	for _, marketEntry := range entry.MarketEntries {
		if result.TotalCost > marketEntry.TotalCost {
			result = marketEntry
		}
	}

	return result
}
