package profitCalc

import (
	csvType "github.com/level-5-pidgey/MarketMoogleApi/csv/types"
	"math"
)

type GcSealExchange struct {
	TokenExchange
	RankRequired int
}

func (gcSealExchange GcSealExchange) GetObtainType() string {
	return "Grand Company Seal"
}

func (gcSealExchange GcSealExchange) GetCost() int {
	return gcSealExchange.Value
}

func (gcSealExchange GcSealExchange) GetQuantity() int {
	return gcSealExchange.Quantity
}

func (gcSealExchange GcSealExchange) GetCostPerItem() int {
	return gcSealExchange.GetCost() / gcSealExchange.GetQuantity()
}

func (gcSealExchange GcSealExchange) GetEffortFactor() float64 {
	return 0.9
}

func calculateSealValue(item *csvType.Item) int {
	if item.Rarity <= 1 || item.ItemLevel == 0 {
		return 0
	}

	const (
		belowIl200Multi        = 5.75
		betweenIl200Il400Multi = 2.00
		betweenIl400Il530Multi = 1.75
		betweenIl530Il660Multi = 1.6667

		betweenIl200Il400Base = 750.0
		betweenIl400Il530Base = 850.50
		betweenIl530Il660Base = 895.0
	)

	if item.ItemLevel <= 200 {
		return int(math.Ceil(float64(item.ItemLevel) * belowIl200Multi))
	} else if item.ItemLevel <= 400 {
		return int(math.Ceil((float64(item.ItemLevel) * betweenIl200Il400Multi) + betweenIl200Il400Base))
	} else if item.ItemLevel <= 530 {
		return int(math.Ceil((float64(item.ItemLevel) * betweenIl400Il530Multi) + betweenIl400Il530Base))
	} else {
		return int(math.Ceil((float64(item.ItemLevel) * betweenIl530Il660Multi) + betweenIl530Il660Base))
	}
}
