package main

import (
	"fmt"
	cache "github.com/go-pkgz/expirable-cache"
	"github.com/level-5-pidgey/MarketMoogle/csv"
	dc "github.com/level-5-pidgey/MarketMoogle/csv/datacollection"
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
	"github.com/level-5-pidgey/MarketMoogle/db"
	"github.com/level-5-pidgey/MarketMoogle/profit"
	"github.com/level-5-pidgey/MarketMoogle/profit/exchange"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	writeWait = 8 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 6) / 10
)

func main() {
	collection, err := dc.CreateDataCollection()
	if err != nil {
		log.Fatal(err)
	}

	profitItems := make(map[int]*profitCalc.Item)
	itemsByObtainInfo := make(map[string]map[int]*profitCalc.Item)
	itemsByExchangeMethod := make(map[string]map[int]*profitCalc.Item)

	for _, csvItem := range *collection.Items {
		item, err := profitCalc.CreateFromCsvData(csvItem, collection)

		if item.ObtainMethods != nil {
			for _, obtainInfo := range *item.ObtainMethods {
				key := obtainInfo.GetExchangeType()

				if itemsByObtainInfo[key] == nil {
					itemsByObtainInfo[key] = make(map[int]*profitCalc.Item)
				}

				itemsByObtainInfo[key][csvItem.Id] = item
			}
		}

		if item.ExchangeMethods != nil {
			for _, exchangeMethod := range *item.ExchangeMethods {
				key := exchangeMethod.GetExchangeType()

				// Don't include dungeon drops to reduce compute time
				if key == readertype.GrandCompanySeal && item.DropsFromDungeon {
					continue
				}

				if itemsByExchangeMethod[key] == nil {
					itemsByExchangeMethod[key] = make(map[int]*profitCalc.Item)
				}

				itemsByExchangeMethod[key][csvItem.Id] = item
			}
		}

		if err != nil {
			log.Fatalf("Error creating item %d: %s", csvItem.Id, err)
		}

		profitItems[csvItem.Id] = item
	}

	// Add Special Shop Currency Exchanges to the profit items map
	for _, shop := range *collection.SpecialShopItem {
		for _, window := range shop.Windows {
			// Don't really want to bother with multi-item exchanges at the moment
			if len(window.Items) > 1 {
				continue
			}

			if len(window.Exchange) > 1 {
				// "Currency" exchanges are denoted with 2 exchanges, with flipped quantities and item ids
				if window.Exchange[0].CostItem != window.Exchange[1].Quantity &&
					window.Exchange[1].CostItem != window.Exchange[0].Quantity {
					continue
				}
			}

			exchangeItem := window.Exchange[0]
			receivedItem := window.Items[0]
			profitItem, ok := profitItems[receivedItem.ItemReceived]

			if !ok {
				continue
			}

			itemCurrency := readertype.FromItemId(exchangeItem.CostItem)

			if itemCurrency == readertype.DefaultCurrency {
				continue
			}

			if profitItem.ObtainMethods == nil {
				obtainMethod := make([]exchange.Method, 0, 1)
				profitItem.ObtainMethods = &obtainMethod
			}

			*profitItem.ObtainMethods = append(
				*profitItem.ObtainMethods, exchange.CurrencyExchange{
					CurrencyType: itemCurrency,
					ShopName:     shop.ShopName,
					Npc:          "", // TODO populate
					Price:        exchangeItem.Quantity,
					Quantity:     receivedItem.Quantity,
				},
			)

			currencyString := itemCurrency.String()
			if itemsByObtainInfo[currencyString] == nil {
				itemsByObtainInfo[currencyString] = make(map[int]*profitCalc.Item)
			}

			itemsByObtainInfo[currencyString][profitItem.Id] = profitItem
		}
	}

	worlds, dataCenters, err := getGameServers()

	if err != nil {
		log.Fatal(err)
	}

	// Create cache
	cacheTime := time.Minute * 10
	c, err := cache.NewCache(cache.TTL(cacheTime))
	if err != nil {
		log.Fatal(err)
	}

	// Periodically empty cache of invalid values
	cacheTicker := time.NewTicker(cacheTime / 2)
	go func(c cache.Cache) {
		for {
			select {
			case <-cacheTicker.C:
				c.DeleteExpired()
			}
		}
	}(c)

	repository, err := db.InitRepository(worlds, dataCenters, c)
	if err != nil {
		log.Fatal(err)
	}

	app := &Application{
		Config: Config{
			Port: os.Getenv("API_PORT"),
		},
	}

	// Create Profit Calculator
	p := profitCalc.NewProfitCalculator(&profitItems, &itemsByObtainInfo, &itemsByExchangeMethod, repository, c)

	// Start up API server
	err = app.Serve(collection, worlds, p)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Started up API server on port %s", app.Config.Port)
}

func getGameServers() (*map[int]*readertype.World, *map[int]*readertype.DataCenter, error) {
	readers := []csv.XivCsvReader{
		csv.UngroupedXivCsvReader[readertype.World]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.World]{
				RowsToSkip: 11,
				FileName:   "World",
			},
		},
		csv.UngroupedXivCsvReader[readertype.DataCenter]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.DataCenter]{
				RowsToSkip: 4,
				FileName:   "WorldDCGroupType",
			},
		},
	}

	var wg sync.WaitGroup
	type csvResults struct {
		data       interface{}
		resultType string
	}

	resultsChan := make(chan csvResults)
	errorsChan := make(chan error)

	for _, reader := range readers {
		wg.Add(1)

		go func(r csv.XivCsvReader) {
			defer wg.Done()

			results, err := r.ProcessCsv()
			if err != nil {
				errorsChan <- err
			} else {
				resultsChan <- csvResults{
					data:       results,
					resultType: r.GetReaderType(),
				}
			}

		}(reader)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
		close(errorsChan)
	}()

	var (
		worlds      map[int]*readertype.World
		dataCenters map[int]*readertype.DataCenter
	)

	results := make([]csvResults, 0)
	readerErrors := make([]error, 0)

	for {
		select {
		case data, ok := <-resultsChan:
			if !ok {
				resultsChan = nil // Avoid reading from closed channel
			} else {
				results = append(results, data)
			}
		case err, ok := <-errorsChan:
			if !ok {
				errorsChan = nil // Avoid reading from closed channel
			} else {
				readerErrors = append(readerErrors, err)
			}
		}

		if resultsChan == nil && errorsChan == nil {
			break // Exit the loop when both channels are closed
		}
	}

	if len(readerErrors) > 0 {
		fmt.Printf("Multiple (%d) readerErrors occurred: ", len(readerErrors))
		for index, err := range readerErrors {
			fmt.Printf("Error #%d: %v\n", index+1, err)
		}

		return nil, nil, fmt.Errorf("multiple (%d) readerErrors occurred", len(readerErrors))
	}

	for _, result := range results {
		switch result.resultType {
		case "World":
			if data, ok := result.data.(map[int]*readertype.World); ok {
				worlds = data
			}
		case "WorldDCGroupType":
			if data, ok := result.data.(map[int]*readertype.DataCenter); ok {
				dataCenters = data
			}
		}
	}

	gameWorlds := make(map[int]*readertype.World)
	gameRegions := map[int]string{
		1: "Japan",
		2: "America",
		3: "Europe",
		4: "Oceania",
	}

	for _, world := range worlds {
		worldDataCenter, ok := dataCenters[world.DataCenterId]
		if ok {
			world.DataCenterName = worldDataCenter.Name
			world.RegionId = worldDataCenter.Group
			world.RegionName = gameRegions[world.RegionId]
		}

		gameWorlds[world.Id] = world
	}

	return &gameWorlds, &dataCenters, nil
}

type Config struct {
	Port string
}

type Application struct {
	Config Config
}

func (app *Application) Serve(
	collection *dc.DataCollection,
	worlds *map[int]*readertype.World,
	profitCalc *profitCalc.ProfitCalculator,
) error {
	port := app.Config.Port

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: Routes(collection, worlds, profitCalc),
	}

	return srv.ListenAndServe()
}
