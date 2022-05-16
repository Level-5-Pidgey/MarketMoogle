/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (xivapiprovider.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package interfaces

import (
	"MarketMoogleAPI/core/apitypes/xivapi"
)

type XivApiProvider interface {
	GetLodestoneInfoById(lodestoneId int) (*xivapi.LodestoneUser, error)
	GetGatheringItemById(contentId int) (*xivapi.GatheringItem, error)
	GetGameItemById(contentId int) (*xivapi.GameItem, error)
	GetRecipeIdByItemId(contentId *int) (*xivapi.RecipeLookup, error)
	GetItems() ([]int, error)
}
