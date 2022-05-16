/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (universalisapiprovider.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package interfaces

import "MarketMoogleAPI/core/apitypes/universalis"

type UniversalisApiProvider interface {
	GetMarketInfoForDc(dataCenter *string, itemId *int) (*universalis.MarketQuery, error)
	GetRecentTransactions(dataCenter *string, amount *int) (*universalis.RecentlyUpdatedItems, error)
	GetMarketableItems() (*[]int, error)
}
