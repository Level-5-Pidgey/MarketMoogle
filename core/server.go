/*
 * Copyright (c) 2022-2023 Carl Alexander Bird.
 * This file (server.go) is part of MarketMoogle and is released under the GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package core

import (
	schema "MarketMoogle/core/graph/model"
	"MarketMoogle/core/util"
	"MarketMoogle/infrastructure/providers/api"
	"MarketMoogle/infrastructure/providers/database"
	"context"
	"fmt"
	"log"
	"time"
)

const defaultPort = "8080"
const initDb = false

func main() {
	log.Println("yippee!")
}

func interval(dbClient *database.Client, servers []string, transCount int) {
	index := 0
	for range time.Tick(time.Minute * 4) {
		if index >= len(servers) {
			index = 0
		}

		err := intervalMarketDataUpdate(dbClient, servers[index], transCount)

		index++

		if err != nil {
			log.Fatal(err)
		}
	}
}

func intervalMarketDataUpdate(dbClient *database.Client, server string, transCount int) error {
	dataCenter := "Materia"

	universalisApiProvider := api.UniversalisApiProvider{}
	marketBoardProvider := database.NewMarketboardDatabaseProvider(dbClient)
	recentTransactions, err := universalisApiProvider.GetRecentTransactions(server, transCount)

	if err != nil {
		log.Print("ran into error getting market info, skipping item")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	for _, item := range recentTransactions.Items {
		marketBoardEntries, err := marketBoardProvider.FindMarketboardEntriesByItemId(ctx, item.ItemID)

		if err != nil {
			log.Fatal(err)
			return err
		}

		for _, marketBoardEntry := range marketBoardEntries {
			if marketBoardEntry.DataCenter != dataCenter {
				continue
			}

			lastUpdateTime, err := util.ConvertTimestampStringToTime(marketBoardEntry.LastUpdateTime)

			if err != nil {
				log.Fatal(err)
				return err
			}

			// If the last update was more than 30 minutes ago, query the API for fresh entries
			if lastUpdateTime.Before(time.Now().UTC().Add(time.Minute * -15)) {
				newMarketData, err := universalisApiProvider.GetMarketInfoForDc(dataCenter, item.ItemID)

				if err != nil {
					log.Print("ran into error getting market info, skipping item")
					return err
				}

				currentTimeString := util.GetCurrentTimestampString()
				err = marketBoardProvider.ReplaceMarketEntry(
					ctx,
					item.ItemID,
					dataCenter,
					newMarketData,
					&currentTimeString,
				)

				if err != nil {
					log.Fatal(err)
					return err
				} else {
					fmt.Printf("%d \t| Updated within %s.\n", item.ItemID, server)
				}
			} else {
				fmt.Printf("%d \t| Skipped updating within %s - already updated recently.\n", item.ItemID, dataCenter)
			}
		}
	}

	return nil
}

func GenerateMarketboardEntries(dbClient *database.Client) {
	universalisApiProvider := api.UniversalisApiProvider{}
	marketBoardProvider := database.NewMarketboardDatabaseProvider(dbClient)
	marketableItems, err := universalisApiProvider.GetMarketableItems()

	if err != nil {
		log.Fatal(err)
	}

	dataCenter := "Materia"
	listingsPerServer := 8
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	for _, itemId := range *marketableItems {
		query, err := universalisApiProvider.GetMarketInfoForDc(dataCenter, itemId)

		if err != nil {
			log.Fatal(err)
		}

		newEntry := schema.MarketboardEntry{
			ItemID:              itemId,
			LastUpdateTime:      util.GetCurrentTimestampString(),
			MarketEntries:       query.GetMarketEntries(listingsPerServer),
			MarketHistory:       query.GetItemHistory(listingsPerServer),
			DataCenter:          query.DcName,
			CurrentAveragePrice: query.CurrentAveragePrice,
			CurrentMinPrice:     &query.MinPrice,
			RegularSaleVelocity: query.RegularSaleVelocity,
			HqSaleVelocity:      query.HqSaleVelocity,
			NqSaleVelocity:      query.NqSaleVelocity,
		}

		returnedEntry, err := marketBoardProvider.CreateMarketEntry(ctx, &newEntry)

		fmt.Printf("%d | ", itemId)
		if returnedEntry != nil {
			fmt.Printf("%d entries on %s added\n", len(returnedEntry.MarketEntries), returnedEntry.DataCenter)
		} else {
			fmt.Printf("0 entries on %s added\n", returnedEntry.DataCenter)
		}
	}
}

func GenerateGameItems(dbClient *database.Client) {
	xivApiProv := api.NewXivApiProvider()
	itemProv := database.NewItemDataBaseProvider(dbClient)
	recipeProv := database.NewRecipeDatabaseProvider(dbClient)

	// Load up all items to add to DB
	allItems, err := xivApiProv.GetItems()

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	// SaveRecipe all items to database
	for _, itemId := range *allItems {
		// Skip 1.0 items that can't be used anymore to save on Db space
		if itemId <= 1601 && itemId >= 19 {
			continue
		}

		gameItem, err := xivApiProv.GetGameItemById(itemId)

		if err != nil {
			log.Fatal(err)
		}

		newItem := schema.Item{
			Id:                 gameItem.ID,
			Name:               gameItem.Name,
			Description:        &gameItem.Description,
			CanBeHq:            gameItem.CanBeHq == 1,
			IconID:             gameItem.IconID,
			SellToVendorValue:  &gameItem.PriceLow,
			BuyFromVendorValue: nil, // This will be added later
		}

		_, err = itemProv.InsertItem(ctx, &newItem)

		// Create recipe (if there is one) and save to DB
		itemRecipes, err := xivApiProv.GetRecipeIdByItemId(itemId)

		if itemRecipes == nil {
			continue
		}

		recipeDict := itemRecipes.GetRecipes()
		for key, value := range recipeDict {
			recipe := value.ConvertToSchemaRecipe(&key)

			_, err = recipeProv.InsertRecipe(ctx, &recipe)

			if err != nil {
				log.Fatal()
			}
		}

		fmt.Printf("%d | %s added", gameItem.ID, gameItem.Name)
		if len(recipeDict) > 0 {
			fmt.Printf(", with %d recipes.", len(recipeDict))
		}
		fmt.Printf("\n")
	}
}

func classNameToEnum(classJob string) schema.CrafterType {
	switch classJob {
	case "CRP":
		return schema.CrafterTypeCarpenter
	case "BSM":
		return schema.CrafterTypeBlacksmith
	case "ARM":
		return schema.CrafterTypeArmourer
	case "GSM":
		return schema.CrafterTypeGoldsmith
	case "LTW":
		return schema.CrafterTypeLeatherworker
	case "WVR":
		return schema.CrafterTypeWeaver
	case "ALC":
		return schema.CrafterTypeAlchemist
	case "CUL":
		return schema.CrafterTypeCulinarian
	default:
		return schema.CrafterType("")
	}
}

func GenerateLeveItems(dbClient *database.Client) {
	xivApiProvider := api.NewXivApiProvider()
	itemProv := database.NewItemDataBaseProvider(dbClient)

	allLeves, err := xivApiProvider.GetCraftLeves()

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	for _, leveId := range *allLeves {
		craftLeve, err := xivApiProvider.GetLeveById(leveId)

		if err != nil {
			leveSchema := schema.CraftingLeve{
				LeveID:    craftLeve.Leve.ID,
				GilReward: craftLeve.Leve.GilReward,
				ExpReward: craftLeve.Leve.ExpReward,
				QuestName: craftLeve.Leve.Name,
				LevelReq:  craftLeve.Leve.ClassJobLevel,
				JobReq:    classNameToEnum(craftLeve.Leve.ClassJobCategory.Name),
			}

			err = itemProv.UpdateLevequestInfo(ctx, craftLeve.Item0.ID, leveSchema)

			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}
}

func GenerateVendorPrices(dbClient *database.Client) {
	xivApiProvider := api.NewXivApiProvider()
	itemProv := database.NewItemDataBaseProvider(dbClient)

	allShops, err := xivApiProvider.GetShops()

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	for _, shopId := range *allShops {
		shopItemsAndPrices, err := xivApiProvider.GetItemsAndPrices(shopId)

		for itemId, price := range shopItemsAndPrices {
			err = itemProv.UpdateVendorSellPrice(ctx, itemId, price)

			if err != nil {
				log.Fatal(err)
				return
			}

			fmt.Printf("%d | Sells for %d and saved to Db\n", itemId, price)
		}
	}
}
