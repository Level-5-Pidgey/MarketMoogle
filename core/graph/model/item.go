/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (item.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package schema

type Item struct {
	ItemID             int     `json:"ItemID"`
	Name               string  `json:"Name"`
	Description        *string `json:"Description"`
	CanBeHq            bool    `json:"CanBeHq"`
	IconID             int     `json:"IconId"`
	SellToVendorValue  *int    `json:"SellToVendorValue"`
	BuyFromVendorValue *int    `json:"BuyFromVendorValue"`
}
