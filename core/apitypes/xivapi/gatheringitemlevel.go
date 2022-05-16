/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (gatheringitemlevel.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package xivapi

type GatheringItemLevel struct {
	GatheringItemLevel int `json:"GatheringItemLevel"`
	ID                 int `json:"ID"`
	Stars              int `json:"Stars"`
}
