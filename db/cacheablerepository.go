package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	cache "github.com/go-pkgz/expirable-cache"
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
	"github.com/level-5-pidgey/MarketMoogle/util"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/pgxpool"
)

const (
	cacheExpiry      = 15 * time.Minute
	ignorePriceValue = 250 * 1000000
	maxRetrievalTime = 172800
	batchLimit       = 100
	listingCacheKey  = "l%d_%v"
	saleCacheKey     = "s%d_%v"
)

type CacheableRepository struct {
	cache cache.Cache

	dataCenters *map[int]*readertype.DataCenter

	worlds *map[int]*readertype.World
}

func makeApiRequest[T any](url string) (*T, error) {
	resp, requestError := http.Get(url)
	if requestError != nil {
		log.Fatal(requestError)
		return nil, requestError
	}

	// API has a DNS problem or is offline, cancel unmarshalling
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("universalis api returned status code %d", resp.StatusCode)
	}

	body, readAllError := io.ReadAll(resp.Body)
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

func InitRepository(
	worlds *map[int]*readertype.World, dataCenters *map[int]*readertype.DataCenter, cache cache.Cache,
) (*CacheableRepository, error) {
	repo := &CacheableRepository{}
	repo.cache = cache
	repo.dataCenters = dataCenters
	repo.worlds = worlds

	return repo, nil
}

func (c *CacheableRepository) GetListingsForItemsOnWorld(itemIds []int, worldId int) (*[]Listing, error) {
	dataCenter := c.getDataCenterFromWorldId(worldId)

	if dataCenter == nil {
		return nil, fmt.Errorf("no data center found for world id %d", worldId)
	}

	listingsOnDc, err := c.GetListingsForItemsOnDataCenter(itemIds, dataCenter.Key)
	if err != nil {
		return nil, err
	}

	var filteredListings []Listing
	for _, listing := range *listingsOnDc {
		if listing.WorldId == worldId {
			filteredListings = append(filteredListings, listing)
		}
	}

	return &filteredListings, err
}

func (c *CacheableRepository) GetListingsForItemsOnDataCenter(itemIds []int, dataCenterId int) (*[]Listing, error) {
	batchedItems := util.BatchItems(itemIds, batchLimit)
	dataCenter := c.getDataCenterFromDcId(dataCenterId)

	if dataCenter == nil {
		return nil, errors.New("data center not found")
	}

	dcString := dataCenter.Name
	var wg sync.WaitGroup
	listingsChan := make(chan []Listing)
	errorsChan := make(chan error, len(batchedItems))

	var results []Listing
	go func() {
		for listings := range listingsChan {
			results = append(results, listings...)
		}
	}()

	for _, batch := range batchedItems {
		wg.Add(1)

		var batchResults []Listing
		go func(batch []int) {
			defer wg.Done()
			// Loop through each batch
			misses := c.getListingCacheMisses(batch, dataCenterId, &batchResults)

			// Get information on all misses from the API
			if len(misses) == 1 {
				listingResults, err := c.updateCacheSingleListing(misses[0], dcString, dataCenterId, dataCenter)
				if err != nil {
					errorsChan <- err
					return
				}

				batchResults = append(batchResults, listingResults...)
			} else if len(misses) > 1 {
				listingResults, err := c.updateCacheMultiListing(misses, dcString, dataCenterId, dataCenter)
				if err != nil {
					errorsChan <- err
					return
				}

				batchResults = append(batchResults, listingResults...)
			}

			listingsChan <- batchResults
		}(batch)
	}

	go func() {
		wg.Wait()
		close(listingsChan)
		close(errorsChan)
	}()

	if len(errorsChan) > 0 {
		fmt.Printf("Multiple (%d) readerErrors occurred: ", len(errorsChan))
		for err := range errorsChan {
			fmt.Printf("Error #%d: %v\n", err)
		}
	}

	return &results, nil
}

func (c *CacheableRepository) updateCacheMultiListing(
	misses []int, dcString string, dataCenterId int, dataCenter *readertype.DataCenter,
) (results []Listing, _ error) {
	missIdsAsString := intArrayToString(misses)
	url := getUniversalisMultiListingUrl(dcString, missIdsAsString)

	res, err := makeApiRequest[MarketData](url)
	if err != nil {
		return nil, err
	}

	// Add apiListings and apiSales to cache
	for itemIdString, marketData := range res.Items {
		itemId, err := strconv.Atoi(itemIdString)
		if err != nil {
			continue
		}

		listings := marketData.Listings
		for _, listing := range listings {
			listing.ItemId = itemId
			listing.DataCenterId = dataCenterId
			listing.RegionId = dataCenter.Group

			results = append(results, listing)
		}

		c.cache.Set(fmt.Sprintf(listingCacheKey, dataCenterId, itemId), listings, cacheExpiry)
	}

	return results, nil
}

func (c *CacheableRepository) updateCacheSingleListing(
	itemId int, dcString string, dataCenterId int, dataCenter *readertype.DataCenter,
) (results []Listing, _ error) {
	url := getUniversalisSingleListingUrl(dcString, itemId)

	res, err := makeApiRequest[ItemDetails](url)
	if err != nil {
		return nil, err
	}

	for _, listing := range res.Listings {
		listing.ItemId = itemId
		listing.DataCenterId = dataCenterId
		listing.RegionId = dataCenter.Group

		results = append(results, listing)
	}

	c.cache.Set(fmt.Sprintf(listingCacheKey, dataCenterId, itemId), results, cacheExpiry)

	return results, nil
}

func (c *CacheableRepository) GetSalesForItemOnDataCenter(itemId, dataCenterId int) (*[]Sale, error) {
	dataCenter := c.getDataCenterFromDcId(dataCenterId)

	if dataCenter == nil {
		return nil, errors.New("data center not found")
	}

	dcString := dataCenter.Name
	cachedSales, found := c.cache.Get(fmt.Sprintf(saleCacheKey, dataCenterId, itemId))
	if found {
		sales := cachedSales.([]Sale)
		return &sales, nil
	}

	foundSales, err := c.updateCacheSingleSale(itemId, dcString, dataCenterId, dataCenter)
	if err != nil {
		return nil, err
	}

	return &foundSales, nil
}

func (c *CacheableRepository) GetSalesForItemsOnWorld(itemIds []int, worldId int) (*[]Sale, error) {
	dataCenter := c.getDataCenterFromWorldId(worldId)

	if dataCenter == nil {
		return nil, fmt.Errorf("no data center found for world id %d", worldId)
	}

	listingsOnDc, err := c.GetSalesForItemsOnDataCenter(itemIds, dataCenter.Key)
	if err != nil {
		return nil, err
	}

	var filteredSales []Sale
	for _, sale := range *listingsOnDc {
		if sale.WorldId == worldId {
			filteredSales = append(filteredSales, sale)
		}
	}

	return &filteredSales, err
}

func (c *CacheableRepository) GetSalesForItemsOnDataCenter(itemIds []int, dataCenterId int) (*[]Sale, error) {
	batchedItems := util.BatchItems(itemIds, batchLimit)
	dataCenter := c.getDataCenterFromDcId(dataCenterId)

	if dataCenter == nil {
		return nil, errors.New("data center not found")
	}

	dcString := dataCenter.Name
	var wg sync.WaitGroup
	salesChan := make(chan []Sale)
	errorsChan := make(chan error, len(batchedItems))

	var results []Sale
	go func() {
		for sales := range salesChan {
			results = append(results, sales...)
		}
	}()

	for _, batch := range batchedItems {
		wg.Add(1)

		var batchResults []Sale
		go func(batch []int) {
			defer wg.Done()

			// Loop through each batch
			misses := c.getSaleCacheMisses(batch, dataCenterId, &batchResults)

			// Get information on all misses from the API
			if len(misses) == 1 {
				foundSales, err := c.updateCacheSingleSale(misses[0], dcString, dataCenterId, dataCenter)
				if err != nil {
					errorsChan <- err
					return
				}

				batchResults = append(batchResults, foundSales...)
			} else if len(misses) > 1 {
				foundSales, err := c.updateCacheMultiSale(misses, dcString, dataCenterId, dataCenter)
				if err != nil {
					errorsChan <- err
					return
				}

				batchResults = append(batchResults, foundSales...)
			}

			salesChan <- batchResults
		}(batch)
	}

	wg.Wait()
	close(salesChan)
	close(errorsChan)

	if len(errorsChan) > 0 {
		fmt.Printf("Multiple (%d) readerErrors occurred: ", len(errorsChan))
		for err := range errorsChan {
			fmt.Printf("Error #%d: %v\n", err)
		}
	}

	return &results, nil
}

func (c *CacheableRepository) updateCacheMultiSale(
	misses []int, dcString string, dataCenterId int, dataCenter *readertype.DataCenter,
) (results []Sale, _ error) {
	missIdsAsString := intArrayToString(misses)
	url := getUniversalisMultiSaleHistoryUrl(dcString, missIdsAsString)

	res, err := makeApiRequest[MarketData](url)
	if err != nil {
		return nil, err
	}

	for itemId, marketData := range res.Items {
		sales := marketData.Sales
		for _, sale := range sales {
			sale.ItemId = marketData.ItemId
			sale.DataCenterId = dataCenterId
			sale.RegionId = dataCenter.Group

			results = append(results, sale)
		}

		c.cache.Set(fmt.Sprintf(saleCacheKey, dataCenterId, itemId), sales, cacheExpiry)
	}

	return results, nil
}

func (c *CacheableRepository) updateCacheSingleSale(
	itemId int, dcString string, dataCenterId int, dataCenter *readertype.DataCenter,
) (results []Sale, _ error) {
	url := getUniversalisSingleSaleHistoryUrl(dcString, itemId)

	res, err := makeApiRequest[ItemDetails](url)
	if err != nil {
		return nil, err
	}

	for _, sale := range res.Sales {
		sale.ItemId = itemId
		sale.DataCenterId = dataCenterId
		sale.RegionId = dataCenter.Group

		results = append(results, sale)
	}

	c.cache.Set(fmt.Sprintf(saleCacheKey, dataCenterId, itemId), results, cacheExpiry)

	return results, nil
}

func (c *CacheableRepository) getSaleCacheMisses(batch []int, dataCenterId int, results *[]Sale) []int {
	misses := make([]int, 0, batchLimit)

	for _, itemId := range batch {
		sales, found := c.cache.Get(fmt.Sprintf(saleCacheKey, dataCenterId, itemId))

		if found {
			cachedSales := sales.([]Sale)
			*results = append(*results, cachedSales...)
		} else {
			misses = append(misses, itemId)
		}
	}

	return misses
}

func (c *CacheableRepository) getListingCacheMisses(batch []int, dataCenterId int, results *[]Listing) []int {
	misses := make([]int, 0, batchLimit)

	for _, itemId := range batch {
		listings, found := c.cache.Get(fmt.Sprintf(listingCacheKey, dataCenterId, itemId))

		if found {
			cachedListings := listings.([]Listing)
			*results = append(*results, cachedListings...)
		} else {
			misses = append(misses, itemId)
		}
	}

	return misses
}

func getUniversalisMultiSaleHistoryUrl(dcString string, missIdsAsString string) string {
	return fmt.Sprintf(
		"https://universalis.app/api/v2/history/%s/%s?entriesToReturn=%d&entriesWithin=%d&maxSalePrice=%d",
		dcString,
		missIdsAsString,
		35,
		maxRetrievalTime,
		ignorePriceValue,
	)
}

func getUniversalisSingleSaleHistoryUrl(dcString string, itemId int) string {
	return fmt.Sprintf(
		"https://universalis.app/api/v2/history/%s/%d?entriesToReturn=%d&entriesWithin=%d&maxSalePrice=%d",
		dcString,
		itemId,
		35,
		maxRetrievalTime,
		ignorePriceValue,
	)
}

func getUniversalisMultiListingUrl(dcString string, missIdsAsString string) string {
	return fmt.Sprintf(
		"https://universalis.app/api/v2/%s/%s?entriesWithin=%d&fields=items.listings%%2Citems.lastUploadTime",
		dcString,
		missIdsAsString,
		maxRetrievalTime,
	)
}

func getUniversalisSingleListingUrl(dcString string, itemId int) string {
	return fmt.Sprintf(
		"https://universalis.app/api/v2/%s/%d?entriesWithin=%d&fields=listings%%2ClastUploadTime",
		dcString,
		itemId,
		maxRetrievalTime,
	)
}

func intArrayToString(i []int) string {
	var buf bytes.Buffer
	for j, v := range i {
		buf.WriteString(strconv.Itoa(v))
		if j < len(i)-1 {
			buf.WriteByte(',')
		}
	}

	return buf.String()
}

func (c *CacheableRepository) getDataCenterFromWorldId(worldId int) *readertype.DataCenter {
	if world, found := (*c.worlds)[worldId]; found {
		return c.getDataCenterFromDcId(world.DataCenterId)
	}

	return nil
}

func (c *CacheableRepository) getDataCenterFromDcId(dcId int) *readertype.DataCenter {
	if dataCenter, found := (*c.dataCenters)[dcId]; found {
		return dataCenter
	}

	return nil
}
