package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/level-5-pidgey/MarketMoogle/api/universalis"
	"github.com/level-5-pidgey/MarketMoogle/csv"
	dc "github.com/level-5-pidgey/MarketMoogle/csv/datacollection"
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
	"github.com/level-5-pidgey/MarketMoogle/db"
	"github.com/level-5-pidgey/MarketMoogle/profit"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

// NA/EU/OCE world Ids
var oceWorldIds = []int{21, 22, 86, 88, 87}
var worldIds = []int{
	21,
	22,
	23,
	24,
	28,
	29,
	30,
	31,
	32,
	33,
	34,
	35,
	36,
	37,
	39,
	40,
	41,
	42,
	43,
	44,
	45,
	46,
	47,
	48,
	49,
	50,
	51,
	52,
	53,
	54,
	55,
	56,
	57,
	58,
	59,
	60,
	61,
	62,
	63,
	64,
	65,
	66,
	67,
	68,
	69,
	70,
	71,
	72,
	73,
	74,
	75,
	76,
	77,
	78,
	79,
	80,
	81,
	82,
	83,
	85,
	86,
	87,
	88,
	90,
	91,
	92,
	93,
	94,
	95,
	96,
	97,
	98,
	99,
	400,
	401,
	402,
	403,
	404,
	405,
	406,
	407,
}

const (
	writeWait = 8 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 6) / 10
)

func main() {
	setupFlag := flag.Bool("setup", false, "runs setup code to initialize db and populate item data")

	flag.Parse()

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
				key := reflect.TypeOf(obtainInfo).String()

				if itemsByObtainInfo[key] == nil {
					itemsByObtainInfo[key] = make(map[int]*profitCalc.Item)
				}

				itemsByObtainInfo[key][csvItem.Id] = item
			}
		}

		if item.ExchangeMethods != nil {
			for _, exchangeMethod := range *item.ExchangeMethods {
				key := reflect.TypeOf(exchangeMethod).String()

				if key == "profitCalc.GcSealExchange" && item.DropsFromDungeon {
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

	worlds, dataCenters, err := getGameServers()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to postgres
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := readDockerSecret("db_user")
	dbPassword := readDockerSecret("db_password")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s timezone=UTC connect_timeout=5",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	repository, err := db.InitRepository(dsn, worlds, dataCenters)
	if err != nil {
		log.Fatal(err)
	}

	defer func(database *pgxpool.Pool) {
		// This doesn't output an anything, so we can't check if
		// there's been an error in the closing process :(
		database.Close()
	}(repository.DbPool)

	app := &Application{
		Config: Config{
			Port: os.Getenv("API_PORT"),
		},
	}

	// Start up API server
	go func() {
		err = app.Serve(&profitItems, &itemsByObtainInfo, &itemsByExchangeMethod, collection, worlds, repository)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Started up API server on port %s", app.Config.Port)
	}()

	// CreatePartitions DB
	if *setupFlag {
		// Create partitions
		err = repository.CreatePartitions()
		if err != nil {
			log.Fatal(err)
		}

		// Get initial listing and sales data with Universalis API
		for _, item := range profitItems {
			if item.MarketProhibited {
				continue
			}

			itemWaitGroup := &sync.WaitGroup{}

			// Search for listings for this item on all data centers
			for _, world := range *worlds {
				itemWaitGroup.Add(1)

				go func(world *readertype.World, group *sync.WaitGroup) {
					defer group.Done()

					listingsUrl := fmt.Sprintf(
						"https://universalis.app/api/v2/%s/%d?listings=40&entries=20",
						strings.ToLower(world.DataCenterName),
						item.Id,
					)

					data, apiErr := makeApiRequest[universalis.Entry](listingsUrl)
					if apiErr != nil && apiErr.Error() != "response object is empty" {
						log.Printf("failed to get listings from universalis: %s\n", apiErr)
					}

					if data == nil {
						return
					}

					// Assign the item id to the data because for some reason universalis uses 2 different "item id" names
					data.Item = item.Id

					listings := data.ConvertToDbListings()

					apiErr = repository.CreateListings(listings)
					if apiErr != nil {
						log.Printf("failed to create listings in db: %s\n", apiErr)
					} else {
						fmt.Printf(
							"Added %d listings for item #%d on the %s datacenter\n",
							len(*listings),
							item.Id,
							world.Name,
						)
					}

					sales := data.ConvertToDbSales()
					apiErr = repository.CreateSales(sales)
					if apiErr != nil {
						log.Printf("failed to create listings in db: %s\n", apiErr)
					} else {
						fmt.Printf(
							"Added %d sales for item #%d on the %s datacenter\n",
							len(*sales),
							item.Id,
							world.Name,
						)
					}
				}(world, itemWaitGroup)
			}

			itemWaitGroup.Wait()
		}
	}

	// Poll Universalis for Market data
	wg := &sync.WaitGroup{}

	for _, worldId := range worldIds {
		wg.Add(1)
		go dialUp(repository, wg, worldId)
	}

	wg.Wait()
}

func makeApiRequest[T any](url string) (*T, error) {
	resp, requestError := http.Get(url)
	if requestError != nil {
		log.Fatal(requestError)
		return nil, requestError
	}

	// API has a DNS problem or is offline, cancel unmarshalling
	if resp.StatusCode == 522 {
		return nil, errors.New("522 code returned from api request")
	}

	body, readAllError := ioutil.ReadAll(resp.Body)
	if readAllError != nil {
		log.Fatal(readAllError)
		return nil, readAllError
	}

	var responseObject T
	var empty T
	err := json.Unmarshal(body, &responseObject)
	if err != nil {
		return nil, err
	}

	// Check if the response object is empty
	if reflect.DeepEqual(responseObject, empty) {
		return nil, errors.New("response object is empty")
	}

	return &responseObject, nil
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

func dialUp(repository db.Repository, wg *sync.WaitGroup, worldId int) {
	defer wg.Done()

	interrupt := make(chan os.Signal, 1)

	u := url.URL{
		Scheme: "wss",
		Host:   "universalis.app",
		Path:   "/api/ws",
	}

	log.Printf("subscribing to ws %s with worldId %d\n", u.String(), worldId)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)

		subscribeToChannel("listings/add", worldId, c)
		subscribeToChannel("sales/add", worldId, c)
		subscribeToChannel("listings/remove", worldId, c)
		subscribeToChannel("sales/remove", worldId, c)

		for {
			msgType, message, err := c.ReadMessage()

			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				break
			}

			if msgType != websocket.BinaryMessage {
				break
			}

			var data universalis.Entry
			err = bson.Unmarshal(message, &data)
			if err != nil {
				log.Printf("failed to unmarshal: %s\n", err)
				return
			}

			switch data.Event {
			case "listings/add":
				dbListings := data.ConvertToDbListings()
				err := repository.CreateListings(dbListings)

				if err != nil {
					log.Printf("failed to create listings in db: %s\n", err)
				}
			case "sales/add":
				dbSales := data.ConvertToDbSales()
				err := repository.CreateSales(dbSales)

				if err != nil {
					log.Printf("failed to create sales in db: %s\n", err)
				}
			case "listings/remove":
				listingIds := make([]string, len(data.Listings))
				for index, listing := range data.Listings {
					listingIds[index] = listing.ListingId
				}

				err := repository.DeleteListings(listingIds)
				if err != nil {
					log.Printf("failed to delete listings in db: %s\n", err)
				}
			case "sales/remove":
				log.Printf("removed sale\n")
			}
		}
	}()

	ticker := time.NewTicker(pingPeriod)

	/*defer func() {
		ticker.Stop()
		log.Printf("closed connection on worldId %d\n", worldId)
		err := c.Close()
		if err != nil {
			return
		}
	}()*/

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err := c.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Println(err)
				return
			}

			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
				return
			}
		case <-interrupt:
			log.Println("interrupted")

			err := c.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)

			if err != nil {
				log.Println("write closed", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func subscribeToChannel(channelName string, worldId int, c *websocket.Conn) {
	subMsg, err := bson.Marshal(
		map[string]string{
			"event":   "subscribe",
			"channel": fmt.Sprintf("%s{world=%d}", channelName, worldId),
		},
	)

	if err != nil {
		log.Fatal("marshal:", err)
	}

	err = c.WriteMessage(websocket.BinaryMessage, subMsg)
	if err != nil {
		log.Fatal("write:", err)
	}
}

func readDockerSecret(secretName string) string {
	secretPath := os.Getenv("SECRETS_DIR") + secretName + os.Getenv("SECRETS_SUFFIX")
	secret, err := os.ReadFile(secretPath)
	if err != nil {
		return ""
	}

	return string(secret)
}

type Config struct {
	Port string
}

type Application struct {
	Config Config
}

func (app *Application) Serve(
	items *map[int]*profitCalc.Item,
	itemsByObtainInfo *map[string]map[int]*profitCalc.Item,
	itemsByExchangeMethod *map[string]map[int]*profitCalc.Item,
	collection *dc.DataCollection,
	worlds *map[int]*readertype.World,
	db db.Repository,
) error {
	port := app.Config.Port

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: Routes(items, itemsByObtainInfo, itemsByExchangeMethod, collection, worlds, db),
	}

	return srv.ListenAndServe()
}
