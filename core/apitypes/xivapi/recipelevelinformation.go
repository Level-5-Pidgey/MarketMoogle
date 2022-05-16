/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (recipelevelinformation.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package xivapi

type RecipeLevelInformation struct {
	ClassJobLevel          int `json:"ClassJobLevel"`
	ConditionsFlag         int `json:"ConditionsFlag"`
	Difficulty             int `json:"Difficulty"`
	Durability             int `json:"Durability"`
	ID                     int `json:"ID"`
	QualityModifier        int `json:"QualityModifier"`
	Stars                  int `json:"Stars"`
	SuggestedControl       int `json:"SuggestedControl"`
	SuggestedCraftsmanship int `json:"SuggestedCraftsmanship"`
}
