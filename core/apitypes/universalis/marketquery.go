/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (marketquery.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package universalis

import (
	schema "MarketMoogle/core/graph/model"
	"sort"
	"strconv"
)

type MarketQuery struct {
	ItemID                int             `json:"itemID"`
	LastUploadTime        int64           `json:"lastUploadTime"`
	Listings              []MarketListing `json:"listings"`
	RecentHistory         []RecentHistory `json:"recentHistory"`
	DcName                string          `json:"dcName"`
	CurrentAveragePrice   float64         `json:"currentAveragePrice"`
	CurrentAveragePriceNQ float64         `json:"currentAveragePriceNQ"`
	CurrentAveragePriceHQ float64         `json:"currentAveragePriceHQ"`
	RegularSaleVelocity   float64         `json:"regularSaleVelocity"`
	NqSaleVelocity        float64         `json:"nqSaleVelocity"`
	HqSaleVelocity        float64         `json:"hqSaleVelocity"`
	AveragePrice          float64         `json:"averagePrice"`
	AveragePriceNQ        float64         `json:"averagePriceNQ"`
	AveragePriceHQ        float64         `json:"averagePriceHQ"`
	MinPrice              int             `json:"minPrice"`
	MinPriceNQ            int             `json:"minPriceNQ"`
	MinPriceHQ            int             `json:"minPriceHQ"`
	MaxPrice              int             `json:"maxPrice"`
	MaxPriceNQ            int             `json:"maxPriceNQ"`
	MaxPriceHQ            int             `json:"maxPriceHQ"`
}

type MarketListing struct {
	LastReviewTime int    `json:"lastReviewTime"`
	PricePerUnit   int    `json:"pricePerUnit"`
	Quantity       int    `json:"quantity"`
	WorldName      string `json:"worldName"`
	WorldID        int    `json:"worldID"`
	Hq             bool   `json:"hq"`
	RetainerName   string `json:"retainerName"`
	Total          int    `json:"total"`
	IsCrafted      bool   `json:"isCrafted"`
}

type RecentHistory struct {
	Hq           bool   `json:"hq"`
	PricePerUnit int    `json:"pricePerUnit"`
	Quantity     int    `json:"quantity"`
	Timestamp    int    `json:"timestamp"`
	WorldName    string `json:"worldName"`
	WorldID      int    `json:"worldID"`
	BuyerName    string `json:"buyerName"`
	Total        int    `json:"total"`
}

func (query MarketQuery) GetItemHistory(listingsPerServer int) []*schema.MarketHistory {
	var marketEntries []*schema.MarketHistory
	servers := make(map[string]int)

	for _, listing := range query.RecentHistory {
		//Only allow a set number of entries per server
		servers[listing.WorldName]++
		if servers[listing.WorldName] >= listingsPerServer {
			continue
		}

		entry := schema.MarketHistory{
			ServerID:        listing.WorldID,
			Server:          listing.WorldName,
			Quantity:        listing.Quantity,
			PricePer:        listing.PricePerUnit,
			TotalPrice:      listing.Total,
			TransactionTime: strconv.Itoa(listing.Timestamp),
			Hq:              listing.Hq,
		}

		marketEntries = append(marketEntries, &entry)
	}

	return marketEntries
}

func (query MarketQuery) GetMarketEntries(listingsPerServer int) []*schema.MarketEntry {
	var marketEntries []*schema.MarketEntry
	servers := make(map[string]int)

	for _, listing := range query.Listings {
		servers[listing.WorldName]++
		if servers[listing.WorldName] >= listingsPerServer {
			continue
		}

		entry := schema.MarketEntry{
			ServerID:     listing.WorldID,
			Server:       listing.WorldName,
			Quantity:     listing.Quantity,
			PricePer:     listing.PricePerUnit,
			TotalCost:    listing.Total,
			Hq:           listing.Hq,
			IsCrafted:    listing.IsCrafted,
			RetainerName: &listing.RetainerName,
		}

		marketEntries = append(marketEntries, &entry)
	}

	//Sort list by price per before returning
	sort.Slice(marketEntries, func(i, j int) bool {
		return marketEntries[i].PricePer < marketEntries[j].PricePer
	})

	return marketEntries
}
