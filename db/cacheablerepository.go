package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	cache "github.com/go-pkgz/expirable-cache"
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/pgxpool"
)

const (
	cacheExpiry      = 15 * time.Minute
	ignorePriceValue = 250 * 1000000
	maxRetrievalTime = 2 * time.Hour * 24 * 365
	batchLimit       = 100
	listingCacheKey  = "l%d_%d"
	saleCacheKey     = "s%d_%d"
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
	if resp.StatusCode == 522 {
		return nil, errors.New("522 code returned from api request")
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

func batchItems(itemIds []int) [][]int {
	var batches [][]int
	for batchLimit < len(itemIds) {
		itemIds, batches = itemIds[batchLimit:], append(batches, itemIds[0:batchLimit:batchLimit])
	}
	batches = append(batches, itemIds)
	return batches
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
	batchedItems := batchItems(itemIds)
	dataCenter := c.getDataCenterFromDcId(dataCenterId)

	if dataCenter == nil {
		return nil, errors.New("data center not found")
	}

	dcString := dataCenter.Name

	var results []Listing
	for _, batch := range batchedItems {
		// Loop through each batch
		misses := make([]int, 0, batchLimit)

		for _, itemId := range batch {
			listing, found := c.cache.Get(fmt.Sprintf(listingCacheKey, dataCenterId, itemId))

			if found {
				cachedListings := listing.([]Listing)
				results = append(results, cachedListings...)
			} else {
				misses = append(misses, itemId)
			}
		}

		// Get information on all misses from the API
		if len(misses) == 1 {
			itemId := misses[0]
			url := fmt.Sprintf(
				"https://universalis.app/api/v2/%s/%d?entriesWithin=%d&fields=listings%%2ClastUploadTime",
				dcString,
				misses[0],
				maxRetrievalTime,
			)

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
		} else if len(misses) > 1 {
			missIdsAsString := intArrayToString(misses)
			url := fmt.Sprintf(
				"https://universalis.app/api/v2/%s/%s?entriesWithin=%d&fields=items.listings%%2Citems.lastUploadTime",
				dcString,
				missIdsAsString,
				maxRetrievalTime,
			)

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

					c.cache.Set(fmt.Sprintf(listingCacheKey, dataCenterId, itemId), listings, cacheExpiry)
					results = append(results, listing)
				}
			}
		}
	}

	return &results, nil
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

	url := fmt.Sprintf(
		"https://universalis.app/api/v2/history/%s/%d?entriesToReturn=%d&entriesWithin=%d&maxSalePrice=%d",
		dcString,
		itemId,
		batchLimit,
		maxRetrievalTime,
		ignorePriceValue,
	)

	res, err := makeApiRequest[ItemDetails](url)
	if err != nil {
		return nil, err
	}

	var sales []Sale
	for _, sale := range res.Sales {
		sale.ItemId = itemId
		sale.RegionId = dataCenter.Group
		sale.DataCenterId = dataCenterId
		sale.Total = sale.Quantity * sale.PricePer

		sales = append(sales, sale)
	}

	c.cache.Set(fmt.Sprintf(saleCacheKey, dataCenterId, itemId), sales, cacheExpiry)

	return &sales, nil
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
	batchedItems := batchItems(itemIds)
	dataCenter := c.getDataCenterFromDcId(dataCenterId)

	if dataCenter == nil {
		return nil, errors.New("data center not found")
	}

	dcString := dataCenter.Name

	var results []Sale
	for _, batch := range batchedItems {
		// Loop through each batch
		misses := make([]int, 0, batchLimit)

		for _, itemId := range batch {
			sales, found := c.cache.Get(fmt.Sprintf(saleCacheKey, dataCenterId, itemId))

			if found {
				cachedSales := sales.([]Sale)
				results = append(results, cachedSales...)
			} else {
				misses = append(misses, itemId)
			}
		}

		// Get information on all misses from the API
		if len(misses) == 1 {
			itemId := misses[0]
			url := fmt.Sprintf(
				"https://universalis.app/api/v2/history/%s/%d?entriesToReturn=%d&entriesWithin=%d&maxSalePrice=%d",
				dcString,
				itemId,
				batchLimit,
				maxRetrievalTime,
				ignorePriceValue,
			)

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
		} else if len(misses) > 0 {
			missIdsAsString := intArrayToString(misses)
			url := fmt.Sprintf(
				"https://universalis.app/api/v2/history/%s/%s?entriesToReturn=%d&entriesWithin=%d&maxSalePrice=%d",
				dcString,
				missIdsAsString,
				batchLimit,
				maxRetrievalTime,
				ignorePriceValue,
			)

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

					c.cache.Set(fmt.Sprintf(saleCacheKey, dataCenterId, itemId), sales, cacheExpiry)
					results = append(results, sale)
				}
			}
		}
	}

	return &results, nil
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
