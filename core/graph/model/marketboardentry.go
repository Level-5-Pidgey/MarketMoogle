/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (marketboardentry.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package schema

import "math"

type MarketboardEntry struct {
	ItemID              int            `json:"ItemId"`
	LastUpdateTime      string         `json:"LastUpdateTime"`
	MarketEntries       []*MarketEntry `json:"MarketEntries"`
	DataCenter          string         `json:"DataCenter"`
	CurrentAveragePrice float64        `json:"CurrentAveragePrice"`
	CurrentMinPrice     *int           `json:"CurrentMinPrice"`
	RegularSaleVelocity float64        `json:"RegularSaleVelocity"`
	HqSaleVelocity      float64        `json:"HqSaleVelocity"`
	NqSaleVelocity      float64        `json:"NqSaleVelocity"`
}

func (m MarketboardEntry) GetCheapestOnServer(server *string) int {
	result := math.MaxInt32

	for _, marketEntry := range m.MarketEntries {
		if marketEntry.Server != *server {
			continue
		}

		if result > marketEntry.PricePer {
			result = marketEntry.PricePer
		}
	}

	return result
}

func (m MarketboardEntry) GetCheapestPriceAndServer() (int, string) {
	resultPrice := math.MaxInt32
	resultServer := ""

	for _, marketEntry := range m.MarketEntries {
		if resultPrice > marketEntry.PricePer {
			resultPrice = marketEntry.PricePer
			resultServer = marketEntry.Server
		}
	}

	return resultPrice, resultServer
}
