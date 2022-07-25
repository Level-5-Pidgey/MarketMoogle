/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (marketboardentry.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package schema

type MarketboardEntry struct {
	ItemID              int              `json:"ItemId"`
	LastUpdateTime      string           `json:"LastUpdateTime"`
	MarketEntries       []*MarketEntry   `json:"MarketEntries"`
	MarketHistory       []*MarketHistory `json:"MarketHistory"`
	DataCenter          string           `json:"DataCenter"`
	CurrentAveragePrice float64          `json:"CurrentAveragePrice"`
	CurrentMinPrice     *int             `json:"CurrentMinPrice"`
	RegularSaleVelocity float64          `json:"RegularSaleVelocity"`
	HqSaleVelocity      float64          `json:"HqSaleVelocity"`
	NqSaleVelocity      float64          `json:"NqSaleVelocity"`
}
