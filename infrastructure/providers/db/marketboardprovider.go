/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (marketboardprovider.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package db

import (
	"MarketMoogleAPI/core/apitypes/universalis"
	schema "MarketMoogleAPI/core/graph/model"
	"MarketMoogleAPI/core/util"
	"MarketMoogleAPI/infrastructure/providers"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math"
	"time"
)

const marketboardCollection = "marketboard"

func (db DbProvider) CreateMarketboardEntryFromApi(dataCenter *string, itemId *int) (*schema.MarketboardEntry, error) {
	//Get different obtain info for item
	prov := providers.UniversalisApiProvider{}

	marketboardListingOut := util.Async(func() *universalis.MarketQuery {
		marketQuery, err := prov.GetMarketInfoForDc(dataCenter, itemId)

		if err != nil {
			log.Fatal(err)
		}

		return marketQuery
	})

	//Turn into item object
	marketListing := <-marketboardListingOut
	marketEntries := marketListing.CreateMarketEntries()

	if len(marketEntries) > 0 {
		marketboardEntry := schema.MarketboardEntry{
			ItemID:              *itemId,
			LastUpdateTime:      util.ConvertTimeToTimestampString(time.Now().UTC()),
			MarketEntries:       marketEntries,
			DataCenter:          marketListing.DcName,
			CurrentAveragePrice: marketListing.CurrentAveragePrice,
			CurrentMinPrice:     &marketListing.MinPrice,
			RegularSaleVelocity: marketListing.RegularSaleVelocity,
			HqSaleVelocity:      marketListing.HqSaleVelocity,
			NqSaleVelocity:      marketListing.NqSaleVelocity,
		}

		return db.SaveMarketboardEntry(&marketboardEntry)
	}

	return nil, nil
}

func (db DbProvider) ReplaceMarketEntries(itemId *int, dataCenter *string, newEntry *universalis.MarketQuery, currentTimestamp *string) error {
	collection := db.client.Database(db.databaseName).Collection(marketboardCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"itemid": *itemId, "datacenter": *dataCenter}

	updatedEntry := schema.MarketboardEntry{
		ItemID:              *itemId,
		LastUpdateTime:      *currentTimestamp,
		MarketEntries:       newEntry.CreateMarketEntries(),
		DataCenter:          *dataCenter,
		CurrentAveragePrice: newEntry.CurrentAveragePrice,
		CurrentMinPrice:     &newEntry.MinPrice,
		RegularSaleVelocity: newEntry.RegularSaleVelocity,
		HqSaleVelocity:      newEntry.HqSaleVelocity,
		NqSaleVelocity:      newEntry.NqSaleVelocity,
	}

	_, err := collection.ReplaceOne(ctx, filter, &updatedEntry, opts)

	return err
}

func (db DbProvider) SaveMarketboardEntry(input *schema.MarketboardEntry) (*schema.MarketboardEntry, error) {
	collection := db.client.Database(db.databaseName).Collection(marketboardCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, input)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return input, nil
}

func (db DbProvider) FindMarketboardEntryByObjectId(ID string) (*schema.MarketboardEntry, error) {
	objectID, err := primitive.ObjectIDFromHex(ID)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	collection := db.client.Database(db.databaseName).Collection(marketboardCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := collection.FindOne(ctx, bson.M{"_id": objectID})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	marketboardEntry := schema.MarketboardEntry{}
	err = result.Decode(&marketboardEntry)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &marketboardEntry, nil
}

func (db DbProvider) FindMarketboardEntriesForDcByItemId(ctx context.Context, itemId int, dataCenter *string) (*schema.MarketboardEntry, error) {
	collection := db.client.Database(db.databaseName).Collection(marketboardCollection)
	cursor, err := collection.Find(ctx, bson.M{"itemid": itemId})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var result schema.MarketboardEntry
	for cursor.Next(ctx) {
		var marketboardEntry schema.MarketboardEntry
		err = cursor.Decode(&marketboardEntry)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		if marketboardEntry.DataCenter == *dataCenter {
			result = marketboardEntry
		}
	}

	return &result, nil
}

func (db DbProvider) FindMarketboardEntriesByItemId(ctx context.Context, itemId int) ([]*schema.MarketboardEntry, error) {
	collection := db.client.Database(db.databaseName).Collection(marketboardCollection)
	cursor, err := collection.Find(ctx, bson.M{"itemid": itemId})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var marketboardEntries []*schema.MarketboardEntry
	for cursor.Next(ctx) {
		var marketboardEntry schema.MarketboardEntry

		err := cursor.Decode(&marketboardEntry)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		marketboardEntries = append(marketboardEntries, &marketboardEntry)
	}

	return marketboardEntries, nil
}

func (db DbProvider) GetAllMarketboardEntries(ctx context.Context) ([]*schema.MarketboardEntry, error) {
	collection := db.client.Database(db.databaseName).Collection(marketboardCollection)
	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var marketboardEntries []*schema.MarketboardEntry
	for cursor.Next(ctx) {
		var item *schema.MarketboardEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		marketboardEntries = append(marketboardEntries, item)
	}

	return marketboardEntries, nil
}

func (db DbProvider) GetCraftingProfit(ctx context.Context, obj *schema.Item, dataCenter *string, homeServer *string, buyCrystals *bool, buyFromOtherServers *bool) (*schema.RecipeResaleInformation, error) {
	var result = schema.RecipeResaleInformation{
		Profit:          0,
		ItemsToPurchase: nil,
		CraftLevel:      0,
	}

	if obj == nil {
		return &result, nil
	}

	itemRecipes, err := db.FindRecipesByItemId(ctx, obj.ItemID)
	if err != nil {
		log.Fatal(err)
		return &result, err
	}

	//If the item has no recipes, there's no crafting profit to be had
	if len(itemRecipes) == 0 {
		return &result, nil
	}

	//Get the sell value of the item on the player's home server
	marketBoardEntriesForItem, err := db.FindMarketboardEntriesForDcByItemId(ctx, obj.ItemID, dataCenter)
	if err != nil {
		log.Fatal(err)
		return &result, err
	}

	if len(marketBoardEntriesForItem.MarketEntries) == 0 {
		return &result, nil
	}

	itemValue := marketBoardEntriesForItem.GetCheapestOnServer(homeServer)
	//If there's no items on the home server, then the value can be assumed from the average price overall
	if itemValue == math.MaxInt32 {
		itemValue = int(marketBoardEntriesForItem.CurrentAveragePrice)
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

				componentItem, err := db.FindItemByItemId(ctx, recipeComponent.ItemID)
				if err != nil {
					log.Fatal(err)
					return &result, err
				}

				itemMarketboardEntries, err := db.FindMarketboardEntriesForDcByItemId(ctx, componentItem.ItemID, dataCenter)
				if err != nil {
					log.Fatal(err)
					return &result, err
				}

				itemPrice, itemServer := getItemPrice(componentItem, itemMarketboardEntries, homeServer, buyFromOtherServers)
				if *itemPrice != math.MaxInt32 {
					itemToPurchase := schema.RecipePurchaseInformation{
						Item:            componentItem,
						ServerToBuyFrom: *itemServer,
						Quantity:        recipeComponent.Count,
					}

					result.ItemsToPurchase = append(result.ItemsToPurchase, &itemToPurchase)

					//Update the cost of this recipe.
					recipeCost += *itemPrice * recipeComponent.Count
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

func getItemPrice(componentItem *schema.Item, itemMarketboardEntries *schema.MarketboardEntry, homeServer *string, buyFromOtherServers *bool) (*int, *string) {
	itemPrice := math.MaxInt32
	purchaseServer := *homeServer

	//Always buy from a vendor if possible, before going to the market board
	if componentItem.BuyFromVendorValue != nil {
		itemPrice = *componentItem.BuyFromVendorValue
	} else {
		cheapestOnMarket := math.MaxInt32
		if buyFromOtherServers != nil && !*buyFromOtherServers { //For lazy people, only check on home server so server swapping is not required
			cheapestOnMarket = itemMarketboardEntries.GetCheapestOnServer(homeServer)
		} else if itemMarketboardEntries.CurrentMinPrice != nil {
			cheapestOnMarket, purchaseServer = itemMarketboardEntries.GetCheapestPriceAndServer()
		}

		if cheapestOnMarket == math.MaxInt32 {
			cheapestOnMarket = int(itemMarketboardEntries.CurrentAveragePrice)
		}

		itemPrice = cheapestOnMarket
	}

	return &itemPrice, &purchaseServer
}

func (db DbProvider) GetVendorResaleProfit(ctx context.Context, obj *schema.Item, dataCenter *string, homeServer *string) (int, error) {
	if obj == nil {
		return 0, nil
	}

	marketEntry, err := db.FindMarketboardEntriesForDcByItemId(ctx, obj.ItemID, dataCenter)

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

	minHomeValue := math.MaxInt32
	for _, entry := range marketEntry.MarketEntries {
		if entry.Server == *homeServer && entry.PricePer < minHomeValue {
			minHomeValue = entry.PricePer
		}
	}

	if minHomeValue == math.MaxInt32 {
		return 0, nil
	}

	return minHomeValue - vendorValue, nil
}

func (db DbProvider) GetCrossDcResaleProfit(ctx context.Context, obj *schema.Item, dataCenter *string, homeServer *string) (int, error) {
	//If no item has been passed through, return early
	if obj == nil {
		return 0, nil
	}

	marketEntry, err := db.FindMarketboardEntriesForDcByItemId(ctx, obj.ItemID, dataCenter)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	//If there are no market entries for this item, return a profit of 0.
	if marketEntry == nil {
		return 0, nil
	}

	cheapestOnHomeServer := math.MaxInt32
	cheapestOnAltServer := math.MaxInt32

	for _, entry := range marketEntry.MarketEntries {
		if entry == nil {
			continue
		}

		if entry.Server == *homeServer {
			if cheapestOnHomeServer > entry.PricePer {
				cheapestOnHomeServer = entry.PricePer
			}
		} else { //If it's not the home server, see if the profit margin is higher
			if cheapestOnAltServer > entry.PricePer {
				cheapestOnAltServer = entry.PricePer
			}
		}
	}

	if cheapestOnAltServer == math.MaxInt32 || cheapestOnHomeServer == math.MaxInt32 {
		return 0, nil
	}

	return cheapestOnHomeServer - cheapestOnAltServer, nil
}
