package exchange

import (
	"bytes"
	"fmt"
	"github.com/level-5-pidgey/MarketMoogle/csv/readertype"
	"math"
)

type GcSealExchange struct {
	CurrencyExchange
	RankRequired readertype.GrandCompanyRank
}

func NewGcSealExchange(price int, npc, shop string, rank readertype.GrandCompanyRank) GcSealExchange {
	return GcSealExchange{
		CurrencyExchange: CurrencyExchange{
			CurrencyType: readertype.GrandCompanySeal,
			Price:        price,
			Quantity:     1,
			Npc:          npc,
			ShopName:     shop,
		},
		RankRequired: rank,
	}
}

func (gcSealExchange GcSealExchange) GetObtainDescription() string {
	var buffer bytes.Buffer
	buffer.WriteString(
		fmt.Sprintf(
			"Exchange %s (Rank: %s",
			gcSealExchange.CurrencyType.GetPlural(),
			gcSealExchange.RankRequired.String(),
		),
	)

	if gcSealExchange.ShopName != "" {
		buffer.WriteString(fmt.Sprintf(", %s)", gcSealExchange.ShopName))
	} else {
		buffer.WriteString(")")
	}

	return buffer.String()
}

func CalculateSealValue(item *readertype.Item) int {
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
