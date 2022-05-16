/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (lodestonecharacter.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package xivapi

type LodestoneCharacter struct {
	Avatar      string      `json:"Avatar"`
	DC          string      `json:"DC"`
	Server      string      `json:"Server"`
	Name        string      `json:"Name"`
	Bio         string      `json:"Bio"`
	ClassJobs   []ClassJob  `json:"ClassJobs"`
	FreeCompany interface{} `json:"FreeCompany"`
}
