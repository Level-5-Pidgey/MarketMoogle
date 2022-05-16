/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (classjob.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package xivapi

type ClassJob struct {
	ClassID       int    `json:"ClassID"`
	ExpLevel      int    `json:"ExpLevel"`
	ExpLevelMax   int    `json:"ExpLevelMax"`
	ExpLevelTogo  int    `json:"ExpLevelTogo"`
	IsSpecialised bool   `json:"IsSpecialised"`
	JobID         int    `json:"JobID"`
	Level         int    `json:"Level"`
	Name          string `json:"Name"`
	UnlockedState struct {
		ID   int    `json:"ID"`
		Name string `json:"Name"`
	} `json:"UnlockedState"`
}
