/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (iteminfo.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package xivapi

type ItemInfo struct {
	AlwaysCollectable          int    `json:"AlwaysCollectable"`
	Article                    int    `json:"Article"`
	CanBeHq                    bool   `json:"CanBeHq"`
	ID                         int    `json:"ID"`
	Icon                       string `json:"Icon"`
	IconHD                     string `json:"IconHD"`
	IconID                     int    `json:"IconID"`
	IsAdvancedMeldingPermitted bool   `json:"IsAdvancedMeldingPermitted"`
	IsCollectable              bool   `json:"IsCollectable"`
	IsCrestWorthy              bool   `json:"IsCrestWorthy"`
	IsDyeable                  bool   `json:"IsDyeable"`
	IsGlamourous               bool   `json:"IsGlamourous"`
	IsIndisposable             bool   `json:"IsIndisposable"`
	IsPvP                      bool   `json:"IsPvP"`
	IsUnique                   bool   `json:"IsUnique"`
	IsUntradable               bool   `json:"IsUntradable"`
	MateriaSlotCount           int    `json:"MateriaSlotCount"`
	Name                       string `json:"Name"`
	Rarity                     int    `json:"Rarity"`
	StackSize                  int    `json:"StackSize"`
	PriceMid                   int    `json:"PriceMid"`
}
