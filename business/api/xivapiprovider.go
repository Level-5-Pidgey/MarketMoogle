/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (xivapiprovider.go) is part of MarketMoogleAPI and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package api

import (
	"MarketMoogle/core/apitypes/xivapi"
)

type XivApiProvider interface {
	GetLodestoneInfoById(lodestoneId int) (*xivapi.LodestoneUser, error)
	GetGatheringItemById(contentId int) (*xivapi.GatheringItem, error)
	GetGameItemById(contentId int) (*xivapi.GameItem, error)
	GetRecipeIdByItemId(contentId *int) (*xivapi.RecipeLookup, error)
	GetItems() ([]int, error)
}
