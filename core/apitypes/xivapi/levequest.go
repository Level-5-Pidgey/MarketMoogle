/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (levequest.go) is part of MarketMoogleAPI and is released under the GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package xivapi

type CraftLeve struct {
	ID   int `json:"ID"`
	Item0 GameItem `json:"Item0"`
	Item0TargetID int `json:"item0TargetID"`
	Leve struct {
		AllowanceCost    int `json:"AllowanceCost"`
		ClassJobCategory struct {
			ID   int    `json:"ID"`
			Name string `json:"Name"`
		} `json:"ClassJobCategory"`
		ClassJobLevel int `json:"ClassJobLevel"`
		DataId        struct {
			ID            int      `json:"ID"`
			Item0         GameItem `json:"Item0"`
			Item0Target   string   `json:"Item0Target"`
			Item0TargetID int      `json:"Item0TargetID"`
		} `json:"DataId"`
		ExpReward       int    `json:"ExpReward"`
		GilReward       int    `json:"GilReward"`
		ID              int    `json:"ID"`
		Name            string `json:"Name"`
		PlaceNameIssued struct {
			ID            int    `json:"ID"`
			Icon          string `json:"Icon"`
			Name          string `json:"Name"`
			NameNoArticle string `json:"NameNoArticle"`
		} `json:"PlaceNameIssued"`
	} `json:"Leve"`
}
