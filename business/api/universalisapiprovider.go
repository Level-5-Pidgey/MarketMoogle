/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (universalisapiprovider.go) is part of MarketMoogleAPI and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package api

import "MarketMoogleAPI/core/apitypes/universalis"

type UniversalisApiProvider interface {
	GetMarketInfoForDc(dataCenter *string, itemId *int) (*universalis.MarketQuery, error)
	GetRecentTransactions(dataCenter *string, amount *int) (*universalis.RecentlyUpdatedItems, error)
	GetMarketableItems() (*[]int, error)
}
