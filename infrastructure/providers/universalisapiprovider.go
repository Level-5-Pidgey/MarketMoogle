/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (universalisapiprovider.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package providers

import (
	"MarketMoogleAPI/core/apitypes/universalis"
	"fmt"
)

type UniversalisApiProvider struct{}

func makeUniversalisContentRequest[T any](contentType string, id *int) (*T, error) {
	url := fmt.Sprintf("https://universalis.app/api/%s/%d", contentType, *id)

	return MakeApiRequest[T](url)
}

func (u UniversalisApiProvider) GetMarketInfoForDc(dataCenter *string, itemId *int) (*universalis.MarketQuery, error) {
	return makeUniversalisContentRequest[universalis.MarketQuery](*dataCenter, itemId)
}

func (u UniversalisApiProvider) GetRecentTransactions(dataCenter *string, amount *int) (*universalis.RecentlyUpdatedItems, error) {
	url := fmt.Sprintf("https://universalis.app/api/extra/stats/most-recently-updated/?dcName=%s&entries=%d", *dataCenter, *amount)

	return MakeApiRequest[universalis.RecentlyUpdatedItems](url)
}

func (u UniversalisApiProvider) GetMarketableItems() (*[]int, error) {
	return MakeApiRequest[[]int]("https://universalis.app/api/marketable")
}
