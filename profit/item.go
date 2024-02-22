package profitCalc

import (
	"errors"
	"github.com/level-5-pidgey/MarketMoogleApi/csv/datacollection"
	csvType "github.com/level-5-pidgey/MarketMoogleApi/csv/types"
)

/*
TODO consolidate types into one domain package. There should only be one "item" type, and when read from the csv
	all this calculation should be done at the same time, or initial information should be layed out and then
	the next readers should add more information to the item
*/

type Item struct {
	Id int
	// Name                 string
	// Description          string
	// IconId               int
	ItemLevel  int
	UiCategory int
	// UiCategoryIconId     int
	// SearchCategory       string
	// SearchCategoryIconId int
	StackSize        int
	Jobs             string
	MarketProhibited bool
	CanBeTraded      bool
	DropsFromDungeon bool
	CanBeHq          bool
	IsCollectable    bool
	// IsGlamour            bool
	ExchangeMethods *[]ExchangeMethod
	ObtainMethods   *[]ExchangeMethod
	CraftingRecipes *[]RecipeInfo
}

func CreateFromCsvData(csvItem *csvType.Item, dataCollection *datacollection.DataCollection) (*Item, error) {
	itemRecipes, recipeError := getRecipes(csvItem, dataCollection)
	obtainMethods, obtainError := getObtainMethods(csvItem, dataCollection)
	exchangeMethods, exchangeError := getExchangeMethods(dataCollection, csvItem)

	if err := errors.Join(recipeError, obtainError, exchangeError); err != nil {
		return nil, err
	}

	// Dereference maps from collection
	uiCategories := *dataCollection.ItemUiCategories
	// searchCategories := *dataCollection.ItemSearchCategories
	classJobCategories := *dataCollection.ClassJobCategories

	uiCatgeory := uiCategories[csvItem.UiCategory]
	// searchCategory := searchCategories[csvItem.SearchCategory]
	jobs := classJobCategories[csvItem.ClassJobCategory].JobCategory

	var result = Item{
		Id: csvItem.Id,
		// Name:                 csvItem.Name,
		// Description:          csvItem.Description,
		// IconId:               csvItem.IconId,
		UiCategory: uiCatgeory.Id,
		// UiCategoryIconId:     uiCatgeory.IconId,
		// SearchCategory:       searchCategory.Name,
		// SearchCategoryIconId: searchCategory.IconId,
		StackSize:        csvItem.StackSize,
		Jobs:             jobs,
		CanBeTraded:      csvItem.CanBeTraded,
		DropsFromDungeon: csvItem.DropsFromDungeon,
		CanBeHq:          csvItem.CanBeHq,
		MarketProhibited: csvItem.SearchCategory == 0,
		IsCollectable:    csvItem.IsCollectable,
		// IsGlamour:            csvItem.IsGlamour,
	}

	// Assign collections if there are any for the item, otherwise leave null.
	if obtainMethods != nil {
		result.ObtainMethods = obtainMethods
	}

	if exchangeMethods != nil {
		result.ExchangeMethods = exchangeMethods
	}

	if itemRecipes != nil {
		result.CraftingRecipes = itemRecipes
	}

	return &result, nil
}

func getExchangeMethods(
	dataCollection *datacollection.DataCollection, csvItem *csvType.Item,
) (*[]ExchangeMethod, error) {
	var exchangeMethods []ExchangeMethod
	if csvItem.SellToVendorPrice > 0 {
		exchangeMethods = append(
			exchangeMethods,
			GilExchange{
				TokenExchange: TokenExchange{
					Value:    csvItem.SellToVendorPrice,
					Quantity: 1,
				},
				NpcName: "NPC",
			},
		)
	}

	if csvItem.Rarity > 1 && csvItem.ItemLevel > 0 {
		exchangeMethods = append(
			exchangeMethods,
			GcSealExchange{
				TokenExchange: TokenExchange{
					Value:    calculateSealValue(csvItem),
					Quantity: 1,
				},
				RankRequired: 6, // Sgt. Second Class
			},
		)
	}

	if len(exchangeMethods) > 0 {
		return &exchangeMethods, nil
	}

	return nil, nil
}
