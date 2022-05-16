/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (content.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package xivapi

type GamePatch struct {
	Banner      string `json:"Banner"`
	ExName      string `json:"ExName"`
	ExVersion   int    `json:"ExVersion"`
	ID          int    `json:"ID"`
	Name        string `json:"Name"`
	ReleaseDate int    `json:"ReleaseDate"`
	Version     string `json:"Version"`
}
