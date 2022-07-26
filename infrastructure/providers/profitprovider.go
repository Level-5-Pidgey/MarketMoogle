/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (profitprovider.go) is part of MarketMoogleAPI and is released under the GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package providers

import (
	schema "MarketMoogleAPI/core/graph/model"
	"MarketMoogleAPI/infrastructure/providers/db"
	"context"
	"log"
	"math"
)

type ItemProfitProvider struct {
	maxValue       int
	returnUnlisted bool

	//Relevant providers
	recipeProvider      *db.RecipeDatabaseProvider
	marketboardProvider *db.MarketboardDatabaseProvider
	itemProvider        *db.ItemDatabaseProvider
}

func NewItemProfitProvider(recipeProv *db.RecipeDatabaseProvider, mbProv *db.MarketboardDatabaseProvider, itemProv *db.ItemDatabaseProvider) *ItemProfitProvider {
	prov := ItemProfitProvider{
		maxValue:            math.MaxInt32,
		returnUnlisted:      true,
		recipeProvider:      recipeProv,
		marketboardProvider: mbProv,
		itemProvider:        itemProv,
	}

	return &prov
}

func (profitProv ItemProfitProvider) GetItemValue(componentItem *schema.Item, mbEntry *schema.MarketboardEntry, homeServer string, buyFromOtherServers *bool) (int, string) {
	itemPrice := profitProv.maxValue
	serverToBuyOn := "vendor"

	//Always buy from a vendor if possible, before going to the market board
	if componentItem.BuyFromVendorValue != nil {
		itemPrice = *componentItem.BuyFromVendorValue
	} else {
		if mbEntry == nil {
			return itemPrice, serverToBuyOn
		}

		cheapestOnMarket := profitProv.maxValue

		//Get cheapest value for the item based on current market entries
		if buyFromOtherServers != nil && !*buyFromOtherServers {
			cheapestOnMarket = profitProv.GetCheapestOnServer(mbEntry, homeServer)
		} else if mbEntry.CurrentMinPrice != nil {
			cheapestOnMarket, homeServer = profitProv.GetCheapestPriceAndServer(mbEntry)
			
		}
		
		//If there are no market entries available, grab the listed average price instead
		if cheapestOnMarket == profitProv.maxValue {
			cheapestOnMarket = int(mbEntry.CurrentAveragePrice)
		}

		itemPrice = cheapestOnMarket
		serverToBuyOn = homeServer
	}

	return itemPrice, serverToBuyOn
}

func (profitProv ItemProfitProvider) GetVendorResaleProfit(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string) (int, error) {
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

func (profitProv ItemProfitProvider) GetCraftingProfit(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string, buyCrystals *bool, buyFromOtherServers *bool) (*schema.RecipeResaleInformation, error) {
	var result = schema.RecipeResaleInformation{
		Profit:          0,
		ItemsToPurchase: nil,
		CraftLevel:      0,
	}

	if obj == nil {
		return &result, nil
	}

	itemRecipes, err := profitProv.recipeProvider.FindRecipesByItemId(ctx, obj.ItemID)
	if err != nil {
		log.Fatal(err)
		return &result, err
	}

	//If the item has no recipes, there's no crafting profit to be had
	if len(itemRecipes) == 0 {
		return &result, nil
	}

	//Get the sell value of the item on the player's home server
	mbEntry, err := profitProv.marketboardProvider.FindItemEntryOnDc(ctx, obj.ItemID, dataCenter)
	if err != nil {
		log.Fatal(err)
		return &result, err
	}

	if mbEntry == nil || len(mbEntry.MarketEntries) == 0 {
		return &result, nil
	}

	itemValue := profitProv.GetCheapestOnServer(mbEntry, homeServer)
	//If there's no items on the home server, then the value can be assumed from the average price overall
	if itemValue == profitProv.maxValue {
		itemValue = int(mbEntry.CurrentAveragePrice)
	}

	recipeCost := 0
	craftType := schema.CraftType("")
	craftLevel := 0

	for i, itemRecipe := range itemRecipes {
		//TODO Check if the player can craft the recipe. Return results for the first recipe.
		if i != 0 {
			continue
		}

		if itemRecipe == nil {
			continue
		}

		craftType = itemRecipe.CraftedBy
		craftLevel = *itemRecipe.RecipeLevel

		recipeComponents := itemRecipe.RecipeItems
		if len(recipeComponents) > 0 {
			for _, recipeComponent := range recipeComponents {
				//If enabled, skip crystals
				if (buyCrystals == nil || !*buyCrystals) && recipeComponent.ItemID <= 19 { //19 is the highest item ID of elemental crystals
					continue
				}

				componentItem, err := profitProv.itemProvider.FindItemByItemId(ctx, recipeComponent.ItemID)
				if err != nil {
					log.Fatal(err)
					return &result, err
				}

				itemMarketboardEntry, err := profitProv.marketboardProvider.FindItemEntryOnDc(ctx, componentItem.ItemID, dataCenter)
				if err != nil {
					log.Fatal(err)
					return &result, err
				}

				itemPrice, itemServer := profitProv.GetItemValue(componentItem, itemMarketboardEntry, homeServer, buyFromOtherServers)
				if itemPrice != profitProv.maxValue {
					itemToPurchase := schema.RecipePurchaseInformation{
						Item:            componentItem,
						ServerToBuyFrom: itemServer,
						Quantity:        recipeComponent.Count,
					}

					result.ItemsToPurchase = append(result.ItemsToPurchase, &itemToPurchase)

					//Update the cost of this recipe.
					recipeCost += itemPrice * recipeComponent.Count
				} else {
					break
				}
			}
		}
	}

	result.ItemCost = recipeCost
	result.Profit = itemValue - recipeCost
	result.CraftType = craftType
	result.CraftLevel = craftLevel

	return &result, nil
}

func (profitProv ItemProfitProvider) GetCrossDcResaleProfit(ctx context.Context, obj *schema.Item, dataCenter string, homeServer string) (int, error) {
	//If no item has been passed through, return early
	if obj == nil {
		return 0, nil
	}

	//Initialize Provider and find item
	marketEntry, err := profitProv.marketboardProvider.FindItemEntryOnDc(ctx, obj.ItemID, dataCenter)

	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	//If there are no market entries for this item, return a profit of 0.
	if marketEntry == nil {
		return 0, nil
	}

	cheapestOnHomeServer := schema.MarketEntry{ PricePer: profitProv.maxValue }
	cheapestOnAltServer := schema.MarketEntry{ PricePer: profitProv.maxValue }
	
	for _, entry := range marketEntry.MarketEntries {
		if entry == nil {
			continue
		}

		if entry.Server == homeServer {
			if cheapestOnHomeServer.PricePer > entry.PricePer {
				cheapestOnHomeServer = *entry
			}
		} else { //If it's not the home server, see if the profit margin is higher
			if cheapestOnAltServer.PricePer > entry.PricePer {
				cheapestOnAltServer = *entry
			}
		}
	}

	if cheapestOnAltServer.PricePer == profitProv.maxValue || cheapestOnHomeServer.PricePer == profitProv.maxValue {
		return 0, nil
	}
	
	return (cheapestOnHomeServer.PricePer * cheapestOnAltServer.Quantity) - cheapestOnAltServer.TotalPrice, nil
}

func (profitProv ItemProfitProvider) GetCheapestOnServer(entry *schema.MarketboardEntry, server string) int {
	result := profitProv.maxValue

	for _, marketEntry := range entry.MarketEntries {
		if marketEntry.Server != server {
			continue
		}

		if result > marketEntry.PricePer {
			result = marketEntry.PricePer
		}
	}

	return result
}

func (profitProv ItemProfitProvider) GetCheapestPriceAndServer(entry *schema.MarketboardEntry) (int, string) {
	resultPrice := profitProv.maxValue
	resultServer := ""

	for _, marketEntry := range entry.MarketEntries {
		if resultPrice > marketEntry.PricePer {
			resultPrice = marketEntry.PricePer
			resultServer = marketEntry.Server
		}
	}

	return resultPrice, resultServer
}
