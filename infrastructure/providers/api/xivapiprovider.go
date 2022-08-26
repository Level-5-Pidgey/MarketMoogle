/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (xivapiprovider.go) is part of MarketMoogleAPI and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogleAPI to view the full license. This file, and all code within MarketMoogleAPI fall under the GNU General Public License.
 */

package api

import (
	"MarketMoogleAPI/core/apitypes/xivapi"
	"fmt"
	"log"
)

type XivApiProvider struct{}

func (p XivApiProvider) GetLodestoneInfoById(lodestoneId *int) (*xivapi.LodestoneUser, error) {
	url := fmt.Sprintf("https://xivapi.com/character/%d", *lodestoneId)

	return MakeApiRequest[xivapi.LodestoneUser](url)
}

func (p XivApiProvider) GetGameItemById(contentId int) (*xivapi.GameItem, error) {
	return MakeXivApiContentRequest[xivapi.GameItem]("Item", contentId)
}

func (p XivApiProvider) GetRecipeIdByItemId(contentId int) (*xivapi.RecipeLookup, error) {
	return MakeXivApiContentRequest[xivapi.RecipeLookup]("RecipeLookup", contentId)
}

func (p XivApiProvider) GetItemsAndPrices(shopId int) (map[int]int, error) {
	gilShop, err := getShopById(shopId)

	if err != nil {
		return nil, err
	}

	itemAndPrice := make(map[int]int)
	for _, item := range gilShop.Items {
		itemAndPrice[item.ID] = item.PriceMid
	}

	return itemAndPrice, nil
}

func getShopById(shopId int) (*xivapi.GilShop, error) {
	return MakeXivApiContentRequest[xivapi.GilShop]("GilShop", shopId)
}

func (p XivApiProvider) GetLeveById(craftLeveId int) (*xivapi.CraftLeve, error) {
	return MakeXivApiContentRequest[xivapi.CraftLeve]("CraftLeve", craftLeveId)
}

func (p XivApiProvider) GetCraftLeves() (*[]int, error) {
	return getPaginatedIds("CraftLeve")
}

func getPaginatedIds(contentType string) (*[]int, error) {
	page := 1
	pageContent, err := MakePaginatedRequest(contentType, page)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	//Create array of items, load from first query
	result := getNonBlankIds(pageContent)

	//Loop through for rest of queries
	for page = 2; page < pageContent.Pagination.PageTotal; page++ {
		fmt.Printf("API Request : Retrieved Page %d\n", page)
		pageContent, err := MakePaginatedRequest(contentType, page)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		result = append(result, getNonBlankIds(pageContent)...)
	}

	return &result, nil
}

func (p XivApiProvider) GetShops() (*[]int, error) {
	return getPaginatedIds("GilShop")
}

func (p XivApiProvider) GetItems() (*[]int, error) {
	return getPaginatedIds("Item")
}

func getNonBlankIds(page *xivapi.PaginatedContent) []int {
	var result []int

	for _, resultContent := range page.Results {
		if resultContent.Name == "" {
			continue
		}

		result = append(result, resultContent.ID)
	}

	return result
}
