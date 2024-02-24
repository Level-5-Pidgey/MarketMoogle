package datacollection

import (
	"fmt"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvReader "github.com/level-5-pidgey/MarketMoogle/csv/readers"
	csvType "github.com/level-5-pidgey/MarketMoogle/csv/types"
	"sync"
)

type DataCollection struct {
	GatheringDataCollection
	RecipeDataCollection
	PlaceDataCollection
	ItemInfoDataCollection
}

func CreateDataCollection() (*DataCollection, error) {
	readers := []csvInterface.GenericCsvReader{
		// Grouped
		csvReader.NewGatheringPointReader(),
		csvReader.NewRecipeCsvReader(),

		// Ungrouped
		csvReader.NewItemCsvReader(),
		csvReader.NewRecipeBookCsvReader(),
		csvReader.NewRecipeLevelReader(),
		csvReader.NewCraftTypeReader(),
		csvReader.NewClassJobCategoryReader(),
		csvReader.NewItemUiCategoryReader(),
		csvReader.NewItemSearchCategoryReader(),
		csvReader.NewGilShopItemReader(),
		csvReader.NewGcScripShopItemReader(),
		csvReader.NewItemGatheringItemReader(),
		csvReader.NewGatheringPointBaseReader(),
		csvReader.NewGatheringItemLevelReader(),
		csvReader.NewItemGatheringTypeReader(),
		csvReader.NewPlaceNameReader(),
		csvReader.NewTerritoryTypeReader(),
	}

	var wg sync.WaitGroup
	type csvResults struct {
		data       interface{}
		resultType string
	}

	resultsChan := make(chan csvResults)
	errorsChan := make(chan error)

	for _, reader := range readers {
		wg.Add(1)

		go func(r csvInterface.GenericCsvReader) {
			defer wg.Done()

			results, err := r.ProcessCsv()
			if err != nil {
				errorsChan <- err
			} else {
				resultsChan <- csvResults{
					data:       results,
					resultType: r.GetReaderType(),
				}
			}

		}(reader)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
		close(errorsChan)
	}()

	var (
		// Grouped
		gatheringPoints map[int][]csvType.GatheringPoint
		recipes         map[int][]csvType.Recipe

		// Ungrouped
		items                map[int]csvType.Item
		recipeBooks          map[int]csvType.RecipeBook
		recipeLevels         map[int]csvType.RecipeLevel
		craftTypes           map[int]csvType.CraftType
		classJobCategories   map[int]csvType.ClassJobCategory
		itemUiCategories     map[int]csvType.ItemUiCategory
		itemSearchCategories map[int]csvType.ItemSearchCategory
		gilShopItems         map[int]csvType.GilShopItem
		gcScripShopItems     map[int]csvType.GcScripShopItem
		gatheringItems       map[int]csvType.GatheringItem
		gatheringPointBases  map[int]csvType.GatheringPointBase
		gatheringItemLevels  map[int]csvType.GatheringItemLevel
		gatheringTypes       map[int]csvType.GatheringType
		placeNames           map[int]csvType.PlaceName
		territoryTypes       map[int]csvType.TerritoryType

		// Misc
		currencies map[int]csvType.Item
	)

	results := make([]csvResults, 0)
	errors := make([]error, 0)

	for {
		select {
		case data, ok := <-resultsChan:
			if !ok {
				resultsChan = nil // Avoid reading from closed channel
			} else {
				results = append(results, data)
			}
		case err, ok := <-errorsChan:
			if !ok {
				errorsChan = nil // Avoid reading from closed channel
			} else {
				errors = append(errors, err)
			}
		}

		if resultsChan == nil && errorsChan == nil {
			break // Exit the loop when both channels are closed
		}
	}

	if len(errors) > 0 {
		fmt.Printf("Multiple (%d) errors occurred: ", len(errors))
		for index, err := range errors {
			fmt.Printf("Error #%d: %v\n", index+1, err)
		}

		return nil, fmt.Errorf("multiple (%d) errors occurred", len(errors))
	}

	dataCollection := DataCollection{
		GatheringDataCollection: GatheringDataCollection{
			GatheringItems:      &gatheringItems,
			GatheringPointBases: &gatheringPointBases,
			GatheringItemLevels: &gatheringItemLevels,
			GatheringTypes:      &gatheringTypes,
			GatheringPoints:     &gatheringPoints,
		},
		RecipeDataCollection: RecipeDataCollection{
			Recipes:      &recipes,
			RecipeBooks:  &recipeBooks,
			RecipeLevels: &recipeLevels,
			CraftTypes:   &craftTypes,
		},
		PlaceDataCollection: PlaceDataCollection{
			PlaceNames:     &placeNames,
			TerritoryTypes: &territoryTypes,
		},
		ItemInfoDataCollection: ItemInfoDataCollection{
			Items:                &items,
			Currencies:           &currencies,
			ClassJobCategories:   &classJobCategories,
			ItemUiCategories:     &itemUiCategories,
			ItemSearchCategories: &itemSearchCategories,
			GilShopItems:         &gilShopItems,
			GcScripShopItem:      &gcScripShopItems,
		},
	}

	for _, result := range results {
		switch result.resultType {
		// Grouped
		case "GatheringPoint":
			if data, ok := result.data.(map[int][]csvType.GatheringPoint); ok {
				gatheringPoints = data
			}
		case "Recipe":
			if data, ok := result.data.(map[int][]csvType.Recipe); ok {
				recipes = data
			}
		// Ungrouped
		case "Item":
			if data, ok := result.data.(map[int]csvType.Item); ok {
				itemMap := make(map[int]csvType.Item)
				currencyMap := make(map[int]csvType.Item)

				// Split currencies into their own map.
				for _, itemData := range data {
					if itemData.SortCategory == 3 ||
						(itemData.SortCategory == 55 && itemData.UiCategory == 61) {
						currencyMap[itemData.Id] = itemData
					} else {
						itemMap[itemData.Id] = itemData
					}
				}

				items = itemMap
				currencies = currencyMap
			}
		case "SecretRecipeBook":
			if data, ok := result.data.(map[int]csvType.RecipeBook); ok {
				recipeBooks = data
			}
		case "RecipeLevelTable":
			if data, ok := result.data.(map[int]csvType.RecipeLevel); ok {
				recipeLevels = data
			}
		case "CraftType":
			if data, ok := result.data.(map[int]csvType.CraftType); ok {
				craftTypes = data
			}
		case "ClassJobCategory":
			if data, ok := result.data.(map[int]csvType.ClassJobCategory); ok {
				classJobCategories = data
			}
		case "ItemUICategory":
			if data, ok := result.data.(map[int]csvType.ItemUiCategory); ok {
				itemUiCategories = data
			}
		case "ItemSearchCategory":
			if data, ok := result.data.(map[int]csvType.ItemSearchCategory); ok {
				itemSearchCategories = data
			}
		case "GilShopItem":
			if data, ok := result.data.(map[int]csvType.GilShopItem); ok {
				gilShopItems = data
			}
		case "GCScripShopItem":
			if data, ok := result.data.(map[int]csvType.GcScripShopItem); ok {
				gcScripShopItems = data
			}
		case "GatheringItem":
			if data, ok := result.data.(map[int]csvType.GatheringItem); ok {
				gatheringItems = data
			}
		case "GatheringPointBase":
			if data, ok := result.data.(map[int]csvType.GatheringPointBase); ok {
				gatheringPointBases = data
			}
		case "GatheringItemLevelConvertTable":
			if data, ok := result.data.(map[int]csvType.GatheringItemLevel); ok {
				gatheringItemLevels = data
			}
		case "GatheringType":
			if data, ok := result.data.(map[int]csvType.GatheringType); ok {
				gatheringTypes = data
			}
		case "PlaceName":
			if data, ok := result.data.(map[int]csvType.PlaceName); ok {
				placeNames = data
			}
		case "TerritoryType":
			if data, ok := result.data.(map[int]csvType.TerritoryType); ok {
				territoryTypes = data
			}
		}
	}

	return &dataCollection, nil
}
