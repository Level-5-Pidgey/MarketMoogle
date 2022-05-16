/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (paginatedcontent.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package xivapi

type PaginatedContent struct {
	Pagination struct {
		Page           int `json:"Page"`
		PageNext       int `json:"PageNext"`
		PagePrev       int `json:"PagePrev"`
		PageTotal      int `json:"PageTotal"`
		Results        int `json:"Results"`
		ResultsPerPage int `json:"ResultsPerPage"`
		ResultsTotal   int `json:"ResultsTotal"`
	} `json:"Pagination"`
	Results []PaginatedResult `json:"Results"`
}

type PaginatedResult struct {
	ID   int         `json:"ID"`
	Icon interface{} `json:"Icon"`
	Name interface{} `json:"Name"`
	Url  string      `json:"Url"`
}
