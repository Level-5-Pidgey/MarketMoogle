/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (server.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package main

import (
	internalGraph "MarketMoogleAPI/core/graph"
	internalGen "MarketMoogleAPI/core/graph/gen"
	schema "MarketMoogleAPI/core/graph/model"
	"MarketMoogleAPI/core/util"
	"MarketMoogleAPI/infrastructure/providers"
	"MarketMoogleAPI/infrastructure/providers/db"
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"time"
)

const defaultPort = "8080"
const initDb = false

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(internalGen.NewExecutableSchema(internalGen.Config{Resolvers: &internalGraph.Resolver{}}))

	//If starting the server for the first time, the mongoDB needs to be populated with items and recipes
	if initDb {
		//TODO Add creation of indexes
		dbProvider := db.NewDbProvider()
		GenerateMarketboardEntries(dbProvider)
		GenerateGameItems(dbProvider)
		GenerateVendorPrices(dbProvider)
	}

	//Ping every 10 minutes
	go interval("Materia", 50)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	//Periodically ping for new market data.

}

func interval(dataCenter string, transCount int) {
	for range time.Tick(time.Minute * 10) {
		err := intervalMarketDataUpdate(&dataCenter, &transCount)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func intervalMarketDataUpdate(dataCenter *string, transCount *int) error {
	universalisApiProvider := providers.UniversalisApiProvider{}
	marketBoardProvider := db.NewDbProvider()
	recentTransactions, err := universalisApiProvider.GetRecentTransactions(dataCenter, transCount)

	if err != nil {
		log.Fatal(err)
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
			if marketBoardEntry.DataCenter != *dataCenter {
				continue
			}

			lastUpdateTime, err := util.ConvertTimestampStringToTime(marketBoardEntry.LastUpdateTime)

			if err != nil {
				log.Fatal(err)
				return err
			}

			//If the last update was more than 30 minutes ago, query the API for fresh entries
			if lastUpdateTime.Before(time.Now().UTC().Add(time.Minute * -15)) {
				newMarketData, err := universalisApiProvider.GetMarketInfoForDc(dataCenter, &item.ItemID)

				if err != nil {
					log.Print("ran into error getting market info, skipping item")
					return err
				}

				currentTimeString := util.GetCurrentTimestampString()
				err = marketBoardProvider.ReplaceMarketEntries(&item.ItemID, dataCenter, newMarketData, &currentTimeString)

				if err != nil {
					log.Fatal(err)
					return err
				} else {
					fmt.Printf("%d \t| Updated within %s.\n", item.ItemID, *dataCenter)
				}
			} else {
				fmt.Printf("%d \t| Skipped updating within %s - already updated recently.\n", item.ItemID, *dataCenter)
			}
		}
	}

	return nil
}

func GenerateMarketboardEntries(dbProv *db.DbProvider) {
	universalisApiProvider := providers.UniversalisApiProvider{}
	marketableItems, err := universalisApiProvider.GetMarketableItems()

	if err != nil {
		log.Fatal(err)
	}

	dataCenter := "Materia"
	for _, itemId := range *marketableItems {
		marketListingOut := util.Async(func() *schema.MarketboardEntry {
			marketListing, err := dbProv.CreateMarketboardEntryFromApi(&dataCenter, &itemId)

			if err != nil {
				log.Fatal(err)
				return nil
			}

			return marketListing
		})

		marketListing := <-marketListingOut

		fmt.Printf("%d | ", itemId)
		if marketListing != nil {
			fmt.Printf("%d entries on %s added\n", len(marketListing.MarketEntries), marketListing.DataCenter)
		} else {
			fmt.Print("0 entries on materia added\n")
		}
	}
}

func GenerateGameItems(dbProv *db.DbProvider) {
	xivApiProvider := providers.XivApiProvider{}

	//Load up all items to add to DB
	allItems, err := xivApiProvider.GetItems()

	if err != nil {
		log.Fatal(err)
	}

	//SaveRecipe all items to database
	for _, itemId := range *allItems {
		//Skip 1.0 items that can't be used anymore to save on Db space
		if itemId <= 1601 && itemId >= 19 {
			continue
		}

		gameItemOut := util.Async(func() *schema.Item {
			gameItem, err := dbProv.SaveItemFromApi(&itemId)

			if err != nil {
				log.Fatal(err)
				return nil
			}

			return gameItem
		})

		//Create recipe (if there is one) and save to DB
		itemRecipeOut := util.Async(func() *[]*schema.Recipe {
			gameRecipes, err := dbProv.CreateRecipesFromApi(&itemId)

			if err != nil {
				log.Fatal(err)
				return nil
			}

			return gameRecipes
		})

		gameItem := <-gameItemOut
		itemRecipes := <-itemRecipeOut

		fmt.Printf("%d | %s added", gameItem.ItemID, gameItem.Name)
		if len(*itemRecipes) > 0 {
			fmt.Printf(", with %d recipes.", len(*itemRecipes))
		}
		fmt.Printf("\n")
	}
}

func GenerateVendorPrices(dbProv *db.DbProvider) {
	xivApiProvider := providers.XivApiProvider{}
	allShops, err := xivApiProvider.GetShops()

	if err != nil {
		log.Fatal(err)
	}

	for _, shopId := range *allShops {
		shopInfoOut := util.Async(func() *map[int]int {
			shopItemsAndPrices, err := xivApiProvider.GetItemsAndPrices(&shopId)

			if err != nil {
				log.Fatal(err)
				return nil
			}

			return shopItemsAndPrices
		})

		shopInfo := <-shopInfoOut

		for itemId, price := range *shopInfo {
			err = dbProv.UpdateVendorSellPrice(&itemId, &price)

			if err != nil {
				log.Fatal(err)
				return
			}

			fmt.Printf("%d | Sells for %d and saved to Db\n", itemId, price)
		}
	}
}
