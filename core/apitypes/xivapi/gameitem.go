/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (gameitem.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package xivapi

type GameItem struct {
	ID                         int    `json:"ID"`
	Url                        string `json:"Url"`
	AlwaysCollectable          int    `json:"AlwaysCollectable"`
	Article                    int    `json:"Article"`
	CanBeHq                    int    `json:"CanBeHq"`
	Icon                       string `json:"Icon"`
	IconHD                     string `json:"IconHD"`
	IconID                     int    `json:"IconID"`
	IsAdvancedMeldingPermitted int    `json:"IsAdvancedMeldingPermitted"`
	IsCollectable              int    `json:"IsCollectable"`
	IsCrestWorthy              int    `json:"IsCrestWorthy"`
	IsDyeable                  int    `json:"IsDyeable"`
	IsGlamourous               int    `json:"IsGlamourous"`
	IsIndisposable             int    `json:"IsIndisposable"`
	IsPvP                      int    `json:"IsPvP"`
	IsUnique                   int    `json:"IsUnique"`
	IsUntradable               int    `json:"IsUntradable"`
	MateriaSlotCount           int    `json:"MateriaSlotCount"`
	Name                       string `json:"Name"`
	Rarity                     int    `json:"Rarity"`
	StackSize                  int    `json:"StackSize"`
	PriceMid                   int    `json:"PriceMid"`
	PriceLow                   int    `json:"PriceLow"`
	Description                string `json:"Description"`
}
