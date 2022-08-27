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
	"MarketMoogleAPI/infrastructure/providers/api"
	"MarketMoogleAPI/infrastructure/providers/database"
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	dbName := "sanctuary"
	var credentials options.Credential
	if os.Getenv("MONGO_DBNAME") != "" {
		dbName = os.Getenv("MONGO_DBNAME")
	}

	//Default credentials if no environment overrides are present
	credentials.AuthSource = "admin"
	credentials.Password = "access123!"
	credentials.Username = "root"
	var hostname = "localhost"
	var mongoPort = "27017"

	if os.Getenv("MONGO_AUTH_SOURCE") != "" {
		credentials.AuthSource = os.Getenv("MONGO_AUTH_SOURCE")
	}

	if os.Getenv("MONGO_PASSWORD") != "" {
		credentials.Password = os.Getenv("MONGO_PASSWORD")
	}

	if os.Getenv("MONGO_USERNAME") != "" {
		credentials.Username = os.Getenv("MONGO_USERNAME")
	}

	if os.Getenv("MONGO_HOST") != "" {
		hostname = os.Getenv("MONGO_HOST")
	}

	if os.Getenv("MONGO_PORT") != "" {
		mongoPort = os.Getenv("MONGO_PORT")
	}

	uri := fmt.Sprintf("mongodb://%s:%s", hostname, mongoPort)

	mongoDbClient := database.NewDatabaseClient(dbName, uri, credentials)

	srv := handler.NewDefaultServer(
		internalGen.NewExecutableSchema(
			internalGen.Config{
				Resolvers: &internalGraph.Resolver{
					DbClient: mongoDbClient,
				},
			}),
	)

	//If starting the server for the first time, the mongoDB needs to be populated with items and recipes
	if initDb {
		CreateDatabaseIndexes(mongoDbClient)
		GenerateMarketboardEntries(mongoDbClient)
		GenerateGameItems(mongoDbClient)
		GenerateVendorPrices(mongoDbClient)
		GenerateLeveItems(mongoDbClient)
	}

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	servers := []string{
		"Ravana",
		"Sophia",
		"Sephirot",
		"Zurvan",
		"Bismarck",
	}
	//Periodically ping for new market data.
	go interval(mongoDbClient, servers, 75)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func interval(dbClient *database.DatabaseClient, servers []string, transCount int) {
	index := 0
	for range time.Tick(time.Minute * 4) {
		if index > len(servers) {
			index = 0
		}

		err := intervalMarketDataUpdate(dbClient, servers[index], transCount)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func intervalMarketDataUpdate(dbClient *database.DatabaseClient, server string, transCount int) error {
	dataCenter := "Materia"

	universalisApiProvider := api.UniversalisApiProvider{}
	marketBoardProvider := database.NewMarketboardDatabaseProvider(dbClient)
	recentTransactions, err := universalisApiProvider.GetRecentTransactions(server, transCount)

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
			if marketBoardEntry.DataCenter != dataCenter {
				continue
			}

			lastUpdateTime, err := util.ConvertTimestampStringToTime(marketBoardEntry.LastUpdateTime)

			if err != nil {
				log.Fatal(err)
				return err
			}

			//If the last update was more than 30 minutes ago, query the API for fresh entries
			if lastUpdateTime.Before(time.Now().UTC().Add(time.Minute * -15)) {
				newMarketData, err := universalisApiProvider.GetMarketInfoForDc(dataCenter, item.ItemID)

				if err != nil {
					log.Print("ran into error getting market info, skipping item")
					return err
				}

				currentTimeString := util.GetCurrentTimestampString()
				err = marketBoardProvider.ReplaceMarketEntry(ctx, item.ItemID, dataCenter, newMarketData, &currentTimeString)

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

func CreateDatabaseIndexes(client *database.DatabaseClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	err := client.CreateIndex(
		ctx,
		"items",
		bson.M{"itemid": 1},
		&options.IndexOptions{Name: util.StringPointer("itemid_index"), Unique: util.BoolPointer(true)},
	)

	if err != nil {
		log.Fatal(err)
	}

	err = client.CreateIndex(
		ctx,
		"marketboard",
		bson.D{{Key: "itemid", Value: 1}, {Key: "datacenter", Value: 1}},
		&options.IndexOptions{Name: util.StringPointer("itemid_and_datacenter_index"), Unique: util.BoolPointer(true)},
	)

	if err != nil {
		log.Fatal(err)
	}

	err = client.CreateIndex(
		ctx,
		"recipes",
		bson.M{"itemresultid": 1},
		&options.IndexOptions{Name: util.StringPointer("itemresult_index"), Unique: util.BoolPointer(false)},
	)

	if err != nil {
		log.Fatal(err)
	}
}

func GenerateMarketboardEntries(dbClient *database.DatabaseClient) {
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

func GenerateGameItems(dbClient *database.DatabaseClient) {
	xivApiProv := api.XivApiProvider{}
	itemProv := database.NewItemDataBaseProvider(dbClient)
	recipeProv := database.NewRecipeDatabaseProvider(dbClient)

	//Load up all items to add to DB
	allItems, err := xivApiProv.GetItems()

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	//SaveRecipe all items to database
	for _, itemId := range *allItems {
		//Skip 1.0 items that can't be used anymore to save on Db space
		if itemId <= 1601 && itemId >= 19 {
			continue
		}

		gameItem, err := xivApiProv.GetGameItemById(itemId)

		if err != nil {
			log.Fatal(err)
		}

		newItem := schema.Item{
			ItemID:             gameItem.ID,
			Name:               gameItem.Name,
			Description:        &gameItem.Description,
			CanBeHq:            gameItem.CanBeHq == 1,
			IconID:             gameItem.IconID,
			SellToVendorValue:  &gameItem.PriceLow,
			BuyFromVendorValue: nil, //This will be added later
		}

		_, err = itemProv.InsertItem(ctx, &newItem)

		//Create recipe (if there is one) and save to DB
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

func GenerateLeveItems(dbClient *database.DatabaseClient) {
	xivApiProvider := api.XivApiProvider{}
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

func GenerateVendorPrices(dbClient *database.DatabaseClient) {
	xivApiProvider := api.XivApiProvider{}
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
