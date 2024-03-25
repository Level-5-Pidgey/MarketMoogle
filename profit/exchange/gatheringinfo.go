package exchange

import (
	"fmt"
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
	"slices"
	"strings"
)

type GatheringInfo struct {
	Points        []GatheringPoint
	Level         int
	Stars         int
	IsCollectible bool
	IsHidden      bool
}

func (gatheringInfo GatheringInfo) GetExchangeType() string {
	return readertype.Gathering
}

func (gatheringInfo GatheringInfo) GetObtainDescription() string {
	var uniqueGatheringTypes []string

	for _, point := range gatheringInfo.Points {
		if !slices.Contains(uniqueGatheringTypes, point.GatherType) {
			uniqueGatheringTypes = append(uniqueGatheringTypes, point.GatherType)
		}
	}

	gatheringTypesString := strings.Join(uniqueGatheringTypes, ", ")
	return fmt.Sprintf("Level %d gathering via %s", gatheringInfo.Level, gatheringTypesString)
}

func (gatheringInfo GatheringInfo) GetEffortFactor() float64 {
	effortFactor := 1.1

	// Add slight penalty for hidden items as they do not always appear at the node
	if gatheringInfo.IsHidden {
		effortFactor += 0.1
	}

	levelFactor := 0.1
	// Reduce effort factor for every 10 levels below the max level this item is to gather
	for i := gatheringInfo.Level; i < 90; i += 10 {
		effortFactor -= levelFactor
		levelFactor /= 2
	}

	return effortFactor
}

func (gatheringInfo GatheringInfo) GetQuantity() int {
	// As a pessimistic guess, you only usually get 3 collectibles per node.
	if gatheringInfo.IsCollectible {
		return 3
	}

	return int(30.0 / gatheringInfo.GetEffortFactor())
}

func (gatheringInfo GatheringInfo) GetCost() int {
	baseCost := 2150

	adjustedCost := baseCost
	for i := gatheringInfo.Level; i < 90; i += 3 {
		// Cap minimum Price of teleport at 150 gil
		if adjustedCost <= 150 {
			return adjustedCost
		}

		adjustedCost -= 75
	}

	return adjustedCost
}

func (gatheringInfo GatheringInfo) GetCostPerItem() float64 {
	costToObtain := gatheringInfo.GetCost()
	amountReceived := gatheringInfo.GetQuantity()

	return float64(costToObtain / amountReceived)
}
