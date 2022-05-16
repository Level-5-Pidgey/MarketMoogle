/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (recentlyupdated.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package universalis

type RecentlyUpdatedItems struct {
	Items []RecentlyUpdatedItem `json:"items"`
}

type RecentlyUpdatedItem struct {
	ItemID         int    `json:"itemID"`
	LastUploadTime int64  `json:"lastUploadTime"`
	WorldID        int    `json:"worldID"`
	WorldName      string `json:"worldName"`
}
