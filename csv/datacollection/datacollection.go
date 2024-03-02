package datacollection

import (
	"fmt"
	"github.com/level-5-pidgey/MarketMoogle/csv"
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
	"sync"
)

type DataCollection struct {
	GatheringDataCollection
	RecipeDataCollection
	PlaceDataCollection
	ItemInfoDataCollection
}

func CreateDataCollection() (*DataCollection, error) {
	readers := []csv.XivCsvReader{
		// Grouped
		csv.GroupedXivCsvReader[readertype.GatheringPoint]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.GatheringPoint]{
				RowsToSkip: 4,
				FileName:   "GatheringPoint",
			},
		},
		csv.GroupedXivCsvReader[readertype.Recipe]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.Recipe]{
				RowsToSkip: 4,
				FileName:   "Recipe",
			},
		},

		// Ungrouped
		csv.UngroupedXivCsvReader[readertype.Item]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.Item]{
				RowsToSkip: 2,
				FileName:   "Item",
			},
		},
		csv.UngroupedXivCsvReader[readertype.RecipeBook]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.RecipeBook]{
				RowsToSkip: 4,
				FileName:   "SecretRecipeBook",
			},
		},
		csv.UngroupedXivCsvReader[readertype.RecipeLevel]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.RecipeLevel]{
				RowsToSkip: 4,
				FileName:   "RecipeLevelTable",
			},
		},
		csv.UngroupedXivCsvReader[readertype.CraftType]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.CraftType]{
				RowsToSkip: 3,
				FileName:   "CraftType",
			},
		},
		csv.UngroupedXivCsvReader[readertype.ClassJobCategory]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.ClassJobCategory]{
				RowsToSkip: 3,
				FileName:   "ClassJobCategory",
			},
		},
		csv.UngroupedXivCsvReader[readertype.ItemUiCategory]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.ItemUiCategory]{
				RowsToSkip: 3,
				FileName:   "ItemUICategory",
			},
		},
		csv.UngroupedXivCsvReader[readertype.ItemSearchCategory]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.ItemSearchCategory]{
				RowsToSkip: 3,
				FileName:   "ItemSearchCategory",
			},
		},
		csv.UngroupedXivCsvReader[readertype.GilShopItem]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.GilShopItem]{
				RowsToSkip: 3,
				FileName:   "GilShopItem",
			},
		},
		csv.UngroupedXivCsvReader[readertype.GcScripShopItem]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.GcScripShopItem]{
				RowsToSkip: 4,
				FileName:   "GCScripShopItem",
			},
		},
		csv.UngroupedXivCsvReader[readertype.GatheringItem]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.GatheringItem]{
				RowsToSkip: 4,
				FileName:   "GatheringItem",
			},
		},
		csv.UngroupedXivCsvReader[readertype.GatheringPointBase]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.GatheringPointBase]{
				RowsToSkip: 4,
				FileName:   "GatheringPointBase",
			},
		},
		csv.UngroupedXivCsvReader[readertype.GatheringItemLevel]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.GatheringItemLevel]{
				RowsToSkip: 4,
				FileName:   "GatheringItemlevelConvertTable",
			},
		},
		csv.UngroupedXivCsvReader[readertype.GatheringType]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.GatheringType]{
				RowsToSkip: 3,
				FileName:   "GatheringType",
			},
		},
		csv.UngroupedXivCsvReader[readertype.PlaceName]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.PlaceName]{
				RowsToSkip: 4,
				FileName:   "PlaceName",
			},
		},
		csv.UngroupedXivCsvReader[readertype.TerritoryType]{
			GenericXivCsvReader: csv.GenericXivCsvReader[readertype.TerritoryType]{
				RowsToSkip: 4,
				FileName:   "TerritoryType",
			},
		},
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

		go func(r csv.XivCsvReader) {
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
		gatheringPoints map[int][]*readertype.GatheringPoint
		recipes         map[int][]*readertype.Recipe

		// Ungrouped
		items                map[int]*readertype.Item
		recipeBooks          map[int]*readertype.RecipeBook
		recipeLevels         map[int]*readertype.RecipeLevel
		craftTypes           map[int]*readertype.CraftType
		classJobCategories   map[int]*readertype.ClassJobCategory
		itemUiCategories     map[int]*readertype.ItemUiCategory
		itemSearchCategories map[int]*readertype.ItemSearchCategory
		gilShopItems         map[int]*readertype.GilShopItem
		gcScripShopItems     map[int]*readertype.GcScripShopItem
		gatheringItems       map[int]*readertype.GatheringItem
		gatheringPointBases  map[int]*readertype.GatheringPointBase
		gatheringItemLevels  map[int]*readertype.GatheringItemLevel
		gatheringTypes       map[int]*readertype.GatheringType
		placeNames           map[int]*readertype.PlaceName
		territoryTypes       map[int]*readertype.TerritoryType

		// Misc
		currencies map[int]*readertype.Item
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
			if data, ok := result.data.(map[int][]*readertype.GatheringPoint); ok {
				gatheringPoints = data
			}
		case "Recipe":
			if data, ok := result.data.(map[int][]*readertype.Recipe); ok {
				recipes = data
			}
		// Ungrouped
		case "Item":
			if data, ok := result.data.(map[int]*readertype.Item); ok {
				itemMap := make(map[int]*readertype.Item)
				currencyMap := make(map[int]*readertype.Item)

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
			if data, ok := result.data.(map[int]*readertype.RecipeBook); ok {
				recipeBooks = data
			}
		case "RecipeLevelTable":
			if data, ok := result.data.(map[int]*readertype.RecipeLevel); ok {
				recipeLevels = data
			}
		case "CraftType":
			if data, ok := result.data.(map[int]*readertype.CraftType); ok {
				craftTypes = data
			}
		case "ClassJobCategory":
			if data, ok := result.data.(map[int]*readertype.ClassJobCategory); ok {
				classJobCategories = data
			}
		case "ItemUICategory":
			if data, ok := result.data.(map[int]*readertype.ItemUiCategory); ok {
				itemUiCategories = data
			}
		case "ItemSearchCategory":
			if data, ok := result.data.(map[int]*readertype.ItemSearchCategory); ok {
				itemSearchCategories = data
			}
		case "GilShopItem":
			if data, ok := result.data.(map[int]*readertype.GilShopItem); ok {
				gilShopItems = data
			}
		case "GCScripShopItem":
			if data, ok := result.data.(map[int]*readertype.GcScripShopItem); ok {
				gcScripShopItems = data
			}
		case "GatheringItem":
			if data, ok := result.data.(map[int]*readertype.GatheringItem); ok {
				gatheringItems = data
			}
		case "GatheringPointBase":
			if data, ok := result.data.(map[int]*readertype.GatheringPointBase); ok {
				gatheringPointBases = data
			}
		case "GatheringItemLevelConvertTable":
			if data, ok := result.data.(map[int]*readertype.GatheringItemLevel); ok {
				gatheringItemLevels = data
			}
		case "GatheringType":
			if data, ok := result.data.(map[int]*readertype.GatheringType); ok {
				gatheringTypes = data
			}
		case "PlaceName":
			if data, ok := result.data.(map[int]*readertype.PlaceName); ok {
				placeNames = data
			}
		case "TerritoryType":
			if data, ok := result.data.(map[int]*readertype.TerritoryType); ok {
				territoryTypes = data
			}
		}
	}

	return &dataCollection, nil
}
