package datacollection

import (
	"fmt"
	csvInterface "github.com/level-5-pidgey/MarketMoogle/csv/interface"
	csvReader "github.com/level-5-pidgey/MarketMoogle/csv/readers"
	"github.com/level-5-pidgey/MarketMoogle/domain"
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
		gatheringPoints map[int][]domain.GatheringPoint
		recipes         map[int][]domain.Recipe

		// Ungrouped
		items                map[int]domain.Item
		recipeBooks          map[int]domain.RecipeBook
		recipeLevels         map[int]domain.RecipeLevel
		craftTypes           map[int]domain.CraftType
		classJobCategories   map[int]domain.ClassJobCategory
		itemUiCategories     map[int]domain.ItemUiCategory
		itemSearchCategories map[int]domain.ItemSearchCategory
		gilShopItems         map[int]domain.GilShopItem
		gcScripShopItems     map[int]domain.GcScripShopItem
		gatheringItems       map[int]domain.GatheringItem
		gatheringPointBases  map[int]domain.GatheringPointBase
		gatheringItemLevels  map[int]domain.GatheringItemLevel
		gatheringTypes       map[int]domain.GatheringType
		placeNames           map[int]domain.PlaceName
		territoryTypes       map[int]domain.TerritoryType

		// Misc
		currencies map[int]domain.Item
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
			if data, ok := result.data.(map[int][]domain.GatheringPoint); ok {
				gatheringPoints = data
			}
		case "Recipe":
			if data, ok := result.data.(map[int][]domain.Recipe); ok {
				recipes = data
			}
		// Ungrouped
		case "Item":
			if data, ok := result.data.(map[int]domain.Item); ok {
				itemMap := make(map[int]domain.Item)
				currencyMap := make(map[int]domain.Item)

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
			if data, ok := result.data.(map[int]domain.RecipeBook); ok {
				recipeBooks = data
			}
		case "RecipeLevelTable":
			if data, ok := result.data.(map[int]domain.RecipeLevel); ok {
				recipeLevels = data
			}
		case "CraftType":
			if data, ok := result.data.(map[int]domain.CraftType); ok {
				craftTypes = data
			}
		case "ClassJobCategory":
			if data, ok := result.data.(map[int]domain.ClassJobCategory); ok {
				classJobCategories = data
			}
		case "ItemUICategory":
			if data, ok := result.data.(map[int]domain.ItemUiCategory); ok {
				itemUiCategories = data
			}
		case "ItemSearchCategory":
			if data, ok := result.data.(map[int]domain.ItemSearchCategory); ok {
				itemSearchCategories = data
			}
		case "GilShopItem":
			if data, ok := result.data.(map[int]domain.GilShopItem); ok {
				gilShopItems = data
			}
		case "GCScripShopItem":
			if data, ok := result.data.(map[int]domain.GcScripShopItem); ok {
				gcScripShopItems = data
			}
		case "GatheringItem":
			if data, ok := result.data.(map[int]domain.GatheringItem); ok {
				gatheringItems = data
			}
		case "GatheringPointBase":
			if data, ok := result.data.(map[int]domain.GatheringPointBase); ok {
				gatheringPointBases = data
			}
		case "GatheringItemLevelConvertTable":
			if data, ok := result.data.(map[int]domain.GatheringItemLevel); ok {
				gatheringItemLevels = data
			}
		case "GatheringType":
			if data, ok := result.data.(map[int]domain.GatheringType); ok {
				gatheringTypes = data
			}
		case "PlaceName":
			if data, ok := result.data.(map[int]domain.PlaceName); ok {
				placeNames = data
			}
		case "TerritoryType":
			if data, ok := result.data.(map[int]domain.TerritoryType); ok {
				territoryTypes = data
			}
		}
	}

	return &dataCollection, nil
}
