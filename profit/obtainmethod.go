package profitCalc

import (
	"github.com/level-5-pidgey/MarketMoogle/csv/datacollection"
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
	"reflect"
	"slices"
)

type ExchangeMethod interface {
	GetObtainType() string

	GetCost() int

	GetQuantity() int

	GetCostPerItem() int

	GetEffortFactor() float64
}

func getObtainMethods(item *readertype.Item, collection *datacollection.DataCollection) (*[]ExchangeMethod, error) {
	scripShopItems := *collection.GcScripShopItem
	gatheringItems := *collection.GatheringItems

	var obtainMethods []ExchangeMethod

	if item.BuyFromVendorPrice > 0 {
		gilShopItems := *collection.GilShopItems

		if _, ok := gilShopItems[item.Id]; ok {
			obtainMethods = append(
				obtainMethods, GilExchange{
					TokenExchange: TokenExchange{
						Value:    item.BuyFromVendorPrice,
						Quantity: 1,
					},
					NpcName: "TBI",
				},
			)
		}
	}

	if gcScripShopItem, ok := scripShopItems[item.Id]; ok {
		obtainMethods = append(
			obtainMethods, GcSealExchange{
				TokenExchange: TokenExchange{
					Value:    gcScripShopItem.AmountRequired,
					Quantity: 1,
				},
				RankRequired: gcScripShopItem.GrandCompanyRankRequired,
			},
		)
	}

	if gatheringItem, ok := gatheringItems[item.Id]; ok {
		obtainMethods = append(
			obtainMethods,
			getGatheringInfo(gatheringItem, collection),
		)
	}

	if len(obtainMethods) > 0 {
		return &obtainMethods, nil
	}

	return nil, nil
}

func getGatheringInfo(
	gatheringItem *readertype.GatheringItem, dataCollection *datacollection.DataCollection,
) *GatheringInfo {
	gatheringLevel := 0
	gatheringStars := 0

	if levelInfo, ok := (*dataCollection.GatheringItemLevels)[gatheringItem.GatheringItemLevelKey]; ok {
		gatheringLevel = levelInfo.Level
		gatheringStars = levelInfo.Stars
	}

	return &GatheringInfo{
		Points:   getGatheringPointsForItem(dataCollection, gatheringItem),
		IsHidden: gatheringItem.IsHidden,
		Level:    gatheringLevel,
		Stars:    gatheringStars,
	}
}

func getGatheringPointsForItem(
	dataCollection *datacollection.DataCollection, gatheringItem *readertype.GatheringItem,
) []GatheringPoint {
	gatheringDataCollection := dataCollection.GatheringDataCollection
	gatheringPointBases := *gatheringDataCollection.GatheringPointBases
	gatheringPoints := *gatheringDataCollection.GatheringPoints
	gatheringTypes := *gatheringDataCollection.GatheringTypes

	gatheringType := ""
	placesToGather := make([]GatheringPoint, 0)

	for _, gatheringPointBase := range gatheringPointBases {
		if slices.Contains(gatheringPointBase.GatheringItemKeys, gatheringItem.Key) {
			gatheringType = gatheringTypes[gatheringPointBase.GatheringTypeKey].Name
			if gPoints, ok := gatheringPoints[gatheringPointBase.Key]; ok {
				if len(gPoints) == 0 {
					continue
				}

				for _, gPoint := range gPoints {
					placeToGather := createPlaceName(
						&dataCollection.PlaceDataCollection,
						gPoint,
						gatheringPointBase,
						gatheringType,
					)

					// Skip points that are hidden and don't have location names
					if placeToGather.Place == "" && placeToGather.Region == "" && placeToGather.Area == "" {
						continue
					}

					// Skip identical gathering points.
					if len(placesToGather) > 0 {
						lastPlace := placesToGather[len(placesToGather)-1]
						if reflect.DeepEqual(lastPlace, placeToGather) {
							continue
						}
					}

					placesToGather = append(placesToGather, placeToGather)
				}
			}
		}
	}

	return placesToGather
}

func createPlaceName(
	placeData *datacollection.PlaceDataCollection,
	gPoint *readertype.GatheringPoint,
	gatheringPointBase *readertype.GatheringPointBase,
	gatheringType string,
) GatheringPoint {
	regionName := ""
	placeName := ""
	areaName := ""

	if pName, ok := (*placeData.PlaceNames)[gPoint.PlaceNameId]; ok {
		areaName = pName.Name
	}

	if territory, ok := (*placeData.TerritoryTypes)[gPoint.TerritoryTypeId]; ok {
		if region, ok := (*placeData.PlaceNames)[territory.RegionId]; ok {
			regionName = region.Name
		}

		if place, ok := (*placeData.PlaceNames)[territory.PlaceId]; ok {
			placeName = place.Name
		}
	}

	return GatheringPoint{
		Level:      gatheringPointBase.GatheringPointLevel,
		GatherType: gatheringType,
		PointType:  "TBI",
		Region:     regionName,
		Place:      placeName,
		Area:       areaName,
	}
}
