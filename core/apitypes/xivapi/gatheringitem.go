/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (gatheringitem.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package xivapi

type GatheringItem struct {
	BaseContent                `json:"ItemInfo"`
	ItemInfo                   ItemInfo           `json:"BaseContent"`
	GatheringItemLevel         GatheringItemLevel `json:"GatheringItemLevel"`
	GatheringItemLevelTargetID int                `json:"GatheringItemLevelTargetID"`
}
