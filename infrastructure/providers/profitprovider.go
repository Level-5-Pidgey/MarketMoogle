/*
 * Copyright (c) 2022 Carl Alexander Bird.
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

func (profitProv ItemProfitProvider) GetRecipePurchaseInfo(componentItem *schema.Item, mbEntry *schema.MarketboardEntry, homeServer string, buyFromOtherServers *bool, count int) *schema.RecipePurchaseInfo {
	result := schema.RecipePurchaseInfo{
		Item:            componentItem,
		ServerToBuyFrom: homeServer,
		BuyFromVendor:   false,
		PricePer:        profitProv.maxValue,
		TotalCost:       profitProv.maxValue,
		Quantity:        count,
	}

	//Always buy from a vendor if possible, before going to the market board
	if componentItem.BuyFromVendorValue != nil {
		result.PricePer = *componentItem.BuyFromVendorValue
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
		if cheapestEntry.TotalCost == profitProv.maxValue {
			result.PricePer = int(mbEntry.CurrentAveragePrice)
		}

		result.PricePer = cheapestEntry.PricePer
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
	if cheapestOnServer.TotalCost == profitProv.maxValue {
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
		PricePer:        profitProv.maxValue,
		TotalCost:       profitProv.maxValue,
		Quantity:        1,
	}

	result := schema.ResaleInfo{
		Profit:          0,
		ItemID:          1,
		Quantity:        1,
		SingleCost:      profitProv.maxValue,
		TotalCost:       profitProv.maxValue,
		ItemsToPurchase: []*schema.RecipePurchaseInfo{},
	}

	//If no item has been passed through, return early
	if obj == nil {
		return nil, nil
	}

	result.ItemID = obj.ItemID

	//Initialize Provider and find item
	marketEntry, err := profitProv.marketboardProvider.FindItemEntryOnDc(ctx, obj.ItemID, dataCenter)
	if err != nil {
		return nil, err
	}

	if marketEntry == nil {
		return nil, errors.New("item searched for is not marketable")
	}

	homeEntry, awayEntry := profitProv.getHomeAndAwayItems(marketEntry, homeServer)

	//Calculate profit per item
	profit := (homeEntry.PricePer * awayEntry.Quantity) - awayEntry.TotalCost

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

	for _, entry := range priceEntries {
		if entry.PricePer < awayEntry.PricePer && entry.PricePer > homeEntry.PricePer {
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

		if marketEntry.PricePer < result.PricePer  {
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