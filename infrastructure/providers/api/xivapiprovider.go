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
	"os"
)

type XivApiProvider struct {
	privateKeyString string
}

func NewXivApiProvider() *XivApiProvider {
	keyString := ""
	if os.Getenv("XIV_API_KEY") != "" {
		keyString = fmt.Sprintf("?private_key=%s", os.Getenv("XIV_API_KEY"))
	}
	
	return &XivApiProvider {
		privateKeyString: keyString,
	}
}
func (p XivApiProvider) GetLodestoneInfoById(lodestoneId int) (*xivapi.LodestoneUser, error) {
	url := fmt.Sprintf("https://xivapi.com/character/%d%s", lodestoneId, p.privateKeyString)

	return MakeApiRequest[xivapi.LodestoneUser](url)
}

func (p XivApiProvider) GetGameItemById(contentId int) (*xivapi.GameItem, error) {
	return MakeXivApiContentRequest[xivapi.GameItem]("Item", contentId, p.privateKeyString)
}

func (p XivApiProvider) GetRecipeIdByItemId(contentId int) (*xivapi.RecipeLookup, error) {
	return MakeXivApiContentRequest[xivapi.RecipeLookup]("RecipeLookup", contentId, p.privateKeyString)
}

func (p XivApiProvider) GetItemsAndPrices(shopId int) (map[int]int, error) {
	gilShop, err := p.getShopById(shopId)

	if err != nil {
		return nil, err
	}

	itemAndPrice := make(map[int]int)
	for _, item := range gilShop.Items {
		itemAndPrice[item.ID] = item.PriceMid
	}

	return itemAndPrice, nil
}

func (p XivApiProvider) getShopById(shopId int) (*xivapi.GilShop, error) {
	return MakeXivApiContentRequest[xivapi.GilShop]("GilShop", shopId, p.privateKeyString)
}

func (p XivApiProvider) GetLeveById(craftLeveId int) (*xivapi.CraftLeve, error) {
	return MakeXivApiContentRequest[xivapi.CraftLeve]("CraftLeve", craftLeveId, p.privateKeyString)
}

func (p XivApiProvider) GetCraftLeves() (*[]int, error) {
	return p.getPaginatedIds("CraftLeve")
}

func (p XivApiProvider) getPaginatedIds(contentType string) (*[]int, error) {
	page := 1
	pageContent, err := MakePaginatedRequest(contentType, page, p.privateKeyString)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	//Create array of items, load from first query
	result := getNonBlankIds(pageContent)

	//Loop through for rest of queries
	for page = 2; page < pageContent.Pagination.PageTotal; page++ {
		fmt.Printf("API Request : Retrieved Page %d\n", page)
		pageContent, err := MakePaginatedRequest(contentType, page, p.privateKeyString)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		result = append(result, getNonBlankIds(pageContent)...)
	}

	return &result, nil
}

func (p XivApiProvider) GetShops() (*[]int, error) {
	return p.getPaginatedIds("GilShop")
}

func (p XivApiProvider) GetItems() (*[]int, error) {
	return p.getPaginatedIds("Item")
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
