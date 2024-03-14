package profitCalc

import (
	"errors"
	"github.com/level-5-pidgey/MarketMoogle/csv/datacollection"
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
	"github.com/level-5-pidgey/MarketMoogle/profit/exchange"
)

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
	Jobs             *[]readertype.Job
	MarketProhibited bool
	CanBeTraded      bool
	DropsFromDungeon bool
	CanBeHq          bool
	IsCollectable    bool
	// IsGlamour            bool
	ExchangeMethods *[]exchange.Method
	ObtainMethods   *[]exchange.Method
	CraftingRecipes *[]RecipeInfo
}

func CreateFromCsvData(csvItem *readertype.Item, dataCollection *datacollection.DataCollection) (*Item, error) {
	itemRecipes, recipeError := getRecipes(csvItem, dataCollection)
	obtainMethods, obtainError := exchange.GetObtainMethods(csvItem, dataCollection)
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

	var jobs *[]readertype.Job
	if category, ok := classJobCategories[csvItem.ClassJobCategory]; ok {
		jobs = &category.JobsInCategory
	}

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
		ItemLevel:        csvItem.ItemLevel,
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
	dataCollection *datacollection.DataCollection, csvItem *readertype.Item,
) (*[]exchange.Method, error) {
	var exchangeMethods []exchange.Method
	if csvItem.SellToVendorPrice > 0 {
		exchangeMethods = append(
			exchangeMethods,
			exchange.NewGilExchange(csvItem.SellToVendorPrice, "", ""), // TODO populate these
		)
	}

	if csvItem.Rarity > 1 &&
		csvItem.ItemLevel > 1 &&
		csvItem.EquipLevel > 1 &&
		csvItem.StackSize == 1 {

		sealPrice := exchange.CalculateSealValue(csvItem)
		exchangeMethods = append(
			exchangeMethods,
			exchange.NewGcSealExchange(sealPrice, "", "", readertype.SergeantSecondClass),
		)
	}

	if len(exchangeMethods) > 0 {
		return &exchangeMethods, nil
	}

	return nil, nil
}
